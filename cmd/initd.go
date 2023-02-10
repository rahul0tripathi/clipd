package main

import (
	"crypto/sha256"
	"fmt"
	"github.com/manifoldco/promptui"
	"github.com/rahul0tripathi/clipd/services/keychain"
	"github.com/rahul0tripathi/clipd/util"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/errors"
)

func notEmpty(s string) error {
	if s == "" {
		return errors.New("value cannot be empty")
	}
	return nil
}
func initd() {
	l, err := util.NewLogger()
	if err != nil {
		panic(err)
	}
	keyServ, err := keychain.NewKeychainService(l)
	if err != nil {
		panic(err)
	}
	fmt.Println("welcome to clipd setup wizard")
	fmt.Println("please enter the master password")
	promptMaster := promptui.Prompt{
		Label:    "Master Password",
		Validate: notEmpty,
		Mask:     '*',
	}
	resultMaster, err := promptMaster.Run()
	if err != nil {
		fmt.Printf("promptMaster failed %v\n", err)
		return
	}
	confirmPromptMaster := promptui.Prompt{
		Label: "Confirm Master Password",
		Mask:  '*',
	}
	confirmPrompt, err := confirmPromptMaster.Run()
	if err != nil {
		fmt.Printf("confirmPromptMaster failed %v\n", err)
		return
	}
	if resultMaster != confirmPrompt {
		fmt.Println("passwords do not match")
		return
	}
	hash := sha256.Sum256([]byte(confirmPrompt))
	err = keyServ.CreateNewSecret(hash[:], nil)
	if err != nil && err != keychain.ErrSecretAlreadyExists {
		fmt.Println("failed to create a new secret in keychain", err)
	}
	if err == keychain.ErrSecretAlreadyExists {
		fmt.Println("a keychain secret with name MasterPassword already exists, please delete before creating new secret")
		fmt.Println("using the existing keychain secret")
	} else {
		fmt.Println("successfully added new password to keychain")
	}
	fmt.Println("please enter the directory where the key permutations will be stored")
	promptKeyDir := promptui.Prompt{
		Label:    "Directory",
		Validate: notEmpty,
	}
	resultKeyDir, err := promptKeyDir.Run()
	if err != nil {
		fmt.Printf("resultKeyDir failed %v\n", err)
		return
	}
	resultKeyDir = fmt.Sprintf("%s/clipd", resultKeyDir)
	allowedPermutations := []int{4, 8, 16, 128, 256, 512}
	promptPermutationNumber := promptui.Select{
		Label: "Select Number Of Permutations To Create",
		Items: allowedPermutations,
	}
	permutationCount, _, err := promptPermutationNumber.Run()
	if err != nil && (permutationCount > len(allowedPermutations) || permutationCount < 0) {
		fmt.Printf("promptPermutationNumber failed %v\n", err)
		return
	}
	db, err := util.NewLevelDB(resultKeyDir, false, l)
	exists, err := db.Get(util.METADATAKEY)
	if err != nil && err != leveldb.ErrNotFound {
		fmt.Printf("db.Get failed %v\n", err)
		return
	}
	if len(exists) != 0 {
		fmt.Println("given directory already has key metadata, please use an empty directory")
		return
	}
	permutations, err := util.GetRandomNPermutations(allowedPermutations[permutationCount])
	if err != nil {
		fmt.Printf("GetRandomNPermutations failed %v\n", err)
		return
	}
	err = util.PersistPermutationToStorage(permutations, db)
	if err != nil {
		fmt.Printf("PersistPermutationToStorage failed %v\n", err)
		return
	}
	fmt.Println("clipd setup wizard completed!")
	fmt.Printf("use the server command with a random salt and the key dir as %s", resultKeyDir)
}
