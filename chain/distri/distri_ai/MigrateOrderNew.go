// Code generated by https://github.com/gagliardetto/anchor-go. DO NOT EDIT.

package distri_ai

import (
	"errors"
	ag_binary "github.com/gagliardetto/binary"
	ag_solanago "github.com/gagliardetto/solana-go"
	ag_format "github.com/gagliardetto/solana-go/text/format"
	ag_treeout "github.com/gagliardetto/treeout"
)

// MigrateOrderNew is the `migrateOrderNew` instruction.
type MigrateOrderNew struct {

	// [0] = [WRITE] orderBefore
	//
	// [1] = [WRITE] orderAfter
	//
	// [2] = [WRITE, SIGNER] signer
	//
	// [3] = [] systemProgram
	ag_solanago.AccountMetaSlice `bin:"-"`
}

// NewMigrateOrderNewInstructionBuilder creates a new `MigrateOrderNew` instruction builder.
func NewMigrateOrderNewInstructionBuilder() *MigrateOrderNew {
	nd := &MigrateOrderNew{
		AccountMetaSlice: make(ag_solanago.AccountMetaSlice, 4),
	}
	return nd
}

// SetOrderBeforeAccount sets the "orderBefore" account.
func (inst *MigrateOrderNew) SetOrderBeforeAccount(orderBefore ag_solanago.PublicKey) *MigrateOrderNew {
	inst.AccountMetaSlice[0] = ag_solanago.Meta(orderBefore).WRITE()
	return inst
}

// GetOrderBeforeAccount gets the "orderBefore" account.
func (inst *MigrateOrderNew) GetOrderBeforeAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(0)
}

// SetOrderAfterAccount sets the "orderAfter" account.
func (inst *MigrateOrderNew) SetOrderAfterAccount(orderAfter ag_solanago.PublicKey) *MigrateOrderNew {
	inst.AccountMetaSlice[1] = ag_solanago.Meta(orderAfter).WRITE()
	return inst
}

// GetOrderAfterAccount gets the "orderAfter" account.
func (inst *MigrateOrderNew) GetOrderAfterAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(1)
}

// SetSignerAccount sets the "signer" account.
func (inst *MigrateOrderNew) SetSignerAccount(signer ag_solanago.PublicKey) *MigrateOrderNew {
	inst.AccountMetaSlice[2] = ag_solanago.Meta(signer).WRITE().SIGNER()
	return inst
}

// GetSignerAccount gets the "signer" account.
func (inst *MigrateOrderNew) GetSignerAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(2)
}

// SetSystemProgramAccount sets the "systemProgram" account.
func (inst *MigrateOrderNew) SetSystemProgramAccount(systemProgram ag_solanago.PublicKey) *MigrateOrderNew {
	inst.AccountMetaSlice[3] = ag_solanago.Meta(systemProgram)
	return inst
}

// GetSystemProgramAccount gets the "systemProgram" account.
func (inst *MigrateOrderNew) GetSystemProgramAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(3)
}

func (inst MigrateOrderNew) Build() *Instruction {
	return &Instruction{BaseVariant: ag_binary.BaseVariant{
		Impl:   inst,
		TypeID: Instruction_MigrateOrderNew,
	}}
}

// ValidateAndBuild validates the instruction parameters and accounts;
// if there is a validation error, it returns the error.
// Otherwise, it builds and returns the instruction.
func (inst MigrateOrderNew) ValidateAndBuild() (*Instruction, error) {
	if err := inst.Validate(); err != nil {
		return nil, err
	}
	return inst.Build(), nil
}

func (inst *MigrateOrderNew) Validate() error {
	// Check whether all (required) accounts are set:
	{
		if inst.AccountMetaSlice[0] == nil {
			return errors.New("accounts.OrderBefore is not set")
		}
		if inst.AccountMetaSlice[1] == nil {
			return errors.New("accounts.OrderAfter is not set")
		}
		if inst.AccountMetaSlice[2] == nil {
			return errors.New("accounts.Signer is not set")
		}
		if inst.AccountMetaSlice[3] == nil {
			return errors.New("accounts.SystemProgram is not set")
		}
	}
	return nil
}

func (inst *MigrateOrderNew) EncodeToTree(parent ag_treeout.Branches) {
	parent.Child(ag_format.Program(ProgramName, ProgramID)).
		//
		ParentFunc(func(programBranch ag_treeout.Branches) {
			programBranch.Child(ag_format.Instruction("MigrateOrderNew")).
				//
				ParentFunc(func(instructionBranch ag_treeout.Branches) {

					// Parameters of the instruction:
					instructionBranch.Child("Params[len=0]").ParentFunc(func(paramsBranch ag_treeout.Branches) {})

					// Accounts of the instruction:
					instructionBranch.Child("Accounts[len=4]").ParentFunc(func(accountsBranch ag_treeout.Branches) {
						accountsBranch.Child(ag_format.Meta("  orderBefore", inst.AccountMetaSlice.Get(0)))
						accountsBranch.Child(ag_format.Meta("   orderAfter", inst.AccountMetaSlice.Get(1)))
						accountsBranch.Child(ag_format.Meta("       signer", inst.AccountMetaSlice.Get(2)))
						accountsBranch.Child(ag_format.Meta("systemProgram", inst.AccountMetaSlice.Get(3)))
					})
				})
		})
}

func (obj MigrateOrderNew) MarshalWithEncoder(encoder *ag_binary.Encoder) (err error) {
	return nil
}
func (obj *MigrateOrderNew) UnmarshalWithDecoder(decoder *ag_binary.Decoder) (err error) {
	return nil
}

// NewMigrateOrderNewInstruction declares a new MigrateOrderNew instruction with the provided parameters and accounts.
func NewMigrateOrderNewInstruction(
	// Accounts:
	orderBefore ag_solanago.PublicKey,
	orderAfter ag_solanago.PublicKey,
	signer ag_solanago.PublicKey,
	systemProgram ag_solanago.PublicKey) *MigrateOrderNew {
	return NewMigrateOrderNewInstructionBuilder().
		SetOrderBeforeAccount(orderBefore).
		SetOrderAfterAccount(orderAfter).
		SetSignerAccount(signer).
		SetSystemProgramAccount(systemProgram)
}
