package cubic

import (
	"testing"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/r1cs"
	"github.com/consensys/gnark/test"
)

func Test_WriteProofInSolidity(t *testing.T) {
	proof, err := ReadProof("./proof")
	if err != nil {
		t.Fatal(err)
	}

	err = WriteProofInSolidity(proof, "./proof_in_sol")
	if err != nil {
		t.Fatal(err)
	}

}

func Test_SolAndProof(t *testing.T) {
	assert := test.NewAssert(t)
	var circuit Circuit

	assignment := &Circuit{
		X: 3,
		Y: 35,
	}

	ccs, err := frontend.Compile(ecc.BN254.ScalarField(), r1cs.NewBuilder, &circuit)
	assert.NoError(err)

	pk, vk, err := groth16.Setup(ccs)
	assert.NoError(err)

	err = WriteVkInSolidity(vk, "./verifier.sol")
	assert.NoError(err)

	witness, err := frontend.NewWitness(assignment, ecc.BN254.ScalarField())
	assert.NoError(err)

	proof, err := groth16.Prove(ccs, pk, witness)
	assert.NoError(err)

	// err = WriteProof(proof, "./proof")
	// assert.NoError(err)

	publicWitness, err := witness.Public()
	assert.NoError(err)

	err = groth16.Verify(proof, vk, publicWitness)
	assert.NoError(err)

	err = WriteProofInSolidity(proof, "./proof_in_sol")
	assert.NoError(err)

	err = WritePublicWitnessInJson(publicWitness, "./public_witness.json")
	assert.NoError(err)
}
