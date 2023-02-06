package util

import (
	"github.com/google/uuid"
	"github.com/syndtr/goleveldb/leveldb/errors"
	"go.uber.org/zap"
	"math/rand"
	"time"
)

const (
	METADATAKEY = "METADATA"
)

type PermutationContext struct {
	Metadata           []string              `json:"metadata"`
	RandomPermutations map[string][256]uint8 `json:"random_permutations"`
}

func GetRandomNPermutations(n int) (*PermutationContext, error) {
	perm := map[string][256]uint8{}
	metadata := []string{}
	for i := 0; i < n; i++ {
		// initialize rand with seed
		rand.Seed(time.Now().UnixNano())
		// initialize the array
		randIndex := uint8(0)
		shuffle := [256]uint8{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61, 62, 63, 64, 65, 66, 67, 68, 69, 70, 71, 72, 73, 74, 75, 76, 77, 78, 79, 80, 81, 82, 83, 84, 85, 86, 87, 88, 89, 90, 91, 92, 93, 94, 95, 96, 97, 98, 99, 100, 101, 102, 103, 104, 105, 106, 107, 108, 109, 110, 111, 112, 113, 114, 115, 116, 117, 118, 119, 120, 121, 122, 123, 124, 125, 126, 127, 128, 129, 130, 131, 132, 133, 134, 135, 136, 137, 138, 139, 140, 141, 142, 143, 144, 145, 146, 147, 148, 149, 150, 151, 152, 153, 154, 155, 156, 157, 158, 159, 160, 161, 162, 163, 164, 165, 166, 167, 168, 169, 170, 171, 172, 173, 174, 175, 176, 177, 178, 179, 180, 181, 182, 183, 184, 185, 186, 187, 188, 189, 190, 191, 192, 193, 194, 195, 196, 197, 198, 199, 200, 201, 202, 203, 204, 205, 206, 207, 208, 209, 210, 211, 212, 213, 214, 215, 216, 217, 218, 219, 220, 221, 222, 223, 224, 225, 226, 227, 228, 229, 230, 231, 232, 233, 234, 235, 236, 237, 238, 239, 240, 241, 242, 243, 244, 245, 246, 247, 248, 249, 250, 251, 252, 253, 254, 255}
		for j := 0; j < 256; j++ {
			randIndex = uint8(rand.Intn(256-j) + j)
			shuffle[j], shuffle[randIndex] = shuffle[randIndex], shuffle[j]
		}
		random, err := uuid.NewRandom()
		if err != nil {
			return nil, err
		}
		// TODO: add check if there is already conflicting key
		perm[random.String()] = shuffle
		metadata = append(metadata, random.String())
	}
	ctx := &PermutationContext{RandomPermutations: perm, Metadata: metadata}
	return ctx, nil
}

func PersistPermutationToStorage(permutations *PermutationContext, storage StorageEngine) error {
	encoded, err := MarshalAndEncode(permutations.Metadata)
	if err != nil {
		return err
	}
	err = storage.Put(METADATAKEY, encoded)
	if err != nil {
		return err
	}
	var encodedPerm []byte
	for uid, perm := range permutations.RandomPermutations {
		encodedPerm, err = MarshalAndEncode(perm)
		if err != nil {
			return err
		}
		err = storage.Put(uid, encodedPerm)
		if err != nil {
			return err
		}
	}
	return nil
}

type PermutedContext struct {
	Key  string
	Data [32]byte
}
type PermutationService interface {
	PermuteRandomly([32]byte) (*PermutedContext, error)
	PermuteWithKey(key string, data [32]byte) ([32]byte, error)
}

type defaultPermutationService struct {
	metadata      []string
	storageEngine StorageEngine
	logger        *zap.Logger
}

func (d defaultPermutationService) PermuteRandomly(data [32]byte) (*PermutedContext, error) {
	rand.Seed(time.Now().UnixNano())
	permuteKey := rand.Intn(len(d.metadata))
	permuted, err := d.PermuteWithKey(d.metadata[permuteKey], data)
	if err != nil {
		return nil, err
	}
	return &PermutedContext{
		Key:  d.metadata[permuteKey],
		Data: permuted,
	}, nil
}

func (d defaultPermutationService) PermuteWithKey(key string, data [32]byte) ([32]byte, error) {
	permutation, err := d.storageEngine.Get(key)
	if err != nil {
		return [32]byte{}, err
	}
	perm := [256]uint8{}
	err = DecodeAndUnmarshal(permutation, &perm)
	if err != nil {
		return [32]byte{}, err
	}
	permuted := [32]byte{}
	currentAt := 0
	counter := 0
	for _, v := range perm {
		mask := byte(1 << (v%8 - 1))
		bit := data[v/8] & mask
		// TODO optimize this
		if bit != 0 {
			permuted[currentAt] |= mask
		} else {
			permuted[currentAt] &^= mask
		}
		counter++
		if counter%8 == 0 {
			currentAt++
			counter = 0
		}
	}
	return permuted, nil
}

func NewDefaultPermutationService(storageService StorageEngine, logger *zap.Logger) (PermutationService, error) {
	service := &defaultPermutationService{metadata: []string{}, logger: logger, storageEngine: storageService}
	metadata, err := storageService.Get(METADATAKEY)

	if err != nil {
		return nil, err
	}
	err = DecodeAndUnmarshal(metadata, &service.metadata)
	if err != nil {
		return nil, err
	}
	if len(service.metadata) == 0 {
		return nil, errors.New("resulting metadata is empty")
	}
	return service, nil
}
