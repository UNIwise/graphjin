package secrets

import (
	"fmt"

	"go.mozilla.org/sops/v3"
	"go.mozilla.org/sops/v3/cmd/sops/codes"
	"go.mozilla.org/sops/v3/cmd/sops/common"
	"go.mozilla.org/sops/v3/keyservice"
)

type decryptOpts struct {
	Cipher      sops.Cipher
	InputStore  sops.Store
	OutputStore sops.Store
	InputPath   string
	IgnoreMAC   bool
	Extract     []interface{}
	KeyServices []keyservice.KeyServiceClient
}

func decrypt(opts decryptOpts) (decryptedFile []byte, err error) {
	tree, err := common.LoadEncryptedFileWithBugFixes(common.GenericDecryptOpts{
		Cipher:      opts.Cipher,
		InputStore:  opts.InputStore,
		InputPath:   opts.InputPath,
		IgnoreMAC:   opts.IgnoreMAC,
		KeyServices: opts.KeyServices,
	})
	if err != nil {
		return nil, err
	}

	_, err = common.DecryptTree(common.DecryptTreeOpts{
		Cipher:      opts.Cipher,
		IgnoreMac:   opts.IgnoreMAC,
		Tree:        tree,
		KeyServices: opts.KeyServices,
	})
	if err != nil {
		return nil, err
	}

	if len(opts.Extract) > 0 {
		return extract(tree, opts.Extract, opts.OutputStore)
	}
	decryptedFile, err = opts.OutputStore.EmitPlainFile(tree.Branches)
	if err != nil {
		return nil, common.NewExitError(fmt.Sprintf("Error dumping file: %s", err), codes.ErrorDumpingTree)
	}
	return decryptedFile, err
}

func extract(tree *sops.Tree, path []interface{}, outputStore sops.Store) (output []byte, err error) {
	v, err := tree.Branches[0].Truncate(path)
	if err != nil {
		return nil, fmt.Errorf("error truncating tree: %s", err)
	}
	if newBranch, ok := v.(sops.TreeBranch); ok {
		tree.Branches[0] = newBranch
		decrypted, err := outputStore.EmitPlainFile(tree.Branches)
		if err != nil {
			return nil, common.NewExitError(fmt.Sprintf("Error dumping file: %s", err), codes.ErrorDumpingTree)
		}
		return decrypted, err
	} else if str, ok := v.(string); ok {
		return []byte(str), nil
	}
	bytes, err := outputStore.EmitValue(v)
	if err != nil {
		return nil, common.NewExitError(fmt.Sprintf("Error dumping tree: %s", err), codes.ErrorDumpingTree)
	}
	return bytes, nil
}
