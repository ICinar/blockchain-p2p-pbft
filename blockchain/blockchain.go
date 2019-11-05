package blockchain


import (
        "time"
        "encoding/hex"
        "golang.org/x/crypto/sha3"

)

type TX struct {
     Timestamp  string
     SignTransm string
     PubKey     string //               from transmitter
     TokenTran  string //string ?       from receiver
}


type Block struct {
        Timestamp string
        Hash      string
        PrevHash  string
        MerkleR TX
}


// make sure block is valid by checking index, and comparing the hash of the previous block
func IsBlockValid(newBlock, oldBlock Block) bool {
        if oldBlock.Hash != newBlock.PrevHash {
                return false
        }

        if CalculateHash(newBlock) != newBlock.Hash {
                return false
        }

        return true
}


// SHA3 hashing
func CalculateHash(block Block) string {
        recordMerk := block.MerkleR.Timestamp + block.MerkleR.SignTransm + block.MerkleR.PubKey + block.MerkleR.TokenTran
        record :=  recordMerk + block.Timestamp + block.PrevHash
        h := sha3.New256()
        h.Write([]byte(record))
        hashed := h.Sum(nil)
        return hex.EncodeToString(hashed)
}

// create a new block using previous block's hash
func GenerateBlock(oldBlock Block, tx TX) Block {

        var newBlock Block

        t := time.Now()
        newBlock.Timestamp = t.String()
        newBlock.MerkleR = tx
        newBlock.PrevHash = oldBlock.Hash
        newBlock.Hash = CalculateHash(newBlock)

        return newBlock
}

