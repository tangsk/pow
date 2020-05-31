
package main
 
import (
"strconv"
"bytes"
"crypto/sha256"
"time"
"fmt"
"math/big"
"encoding/binary"
"log"
"math"
"os"
)
 
const targetBits = 10
var maxNonce = math.MaxInt64
 
type Block struct {
	Timestamp     int64
	Data          []byte
	PrevBlockHash []byte
	Hash          []byte
	Nonce         int

}
 
type ProofOfWork struct {
	block  *Block
	target *big.Int
}
 
func (b *Block) SetHash() {
	timestamp := []byte(strconv.FormatInt(b.Timestamp, 10))
	headers := bytes.Join([][]byte{b.PrevBlockHash, b.Data, timestamp}, []byte{})
	hash := sha256.Sum256(headers)
 
	b.Hash = hash[:]
}
 
func NewProofOfWork(b *Block) *ProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-targetBits))
 
	pow := &ProofOfWork{b, target}
 
	return pow
}
 
func IntToHex(num int64) []byte {
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, num)
	if err != nil {
		log.Panic(err)
	}
 
	return buff.Bytes()
}
 
func (pow *ProofOfWork) prepareData(nonce int) []byte {
	data := bytes.Join(
		[][]byte{
			pow.block.PrevBlockHash,
			pow.block.Data,
			IntToHex(pow.block.Timestamp),
			IntToHex(int64(targetBits)),
			IntToHex(int64(nonce)),
		},
		[]byte{},
	)
 
	return data
}
 
func (pow *ProofOfWork) Run() (int, []byte) {
	var hashInt big.Int
	var hash [32]byte
	nonce := 0
 
	fmt.Printf("Mining the block containing \"%s\"\n", pow.block.Data)
	for nonce < maxNonce {
		data := pow.prepareData(nonce)
		hash = sha256.Sum256(data)
		fmt.Printf("\r%x", hash)
		hashInt.SetBytes(hash[:])
 
		if hashInt.Cmp(pow.target) == -1 {
			break
		} else {
			nonce++
		}
	}
	fmt.Println()
	fmt.Println(nonce)
 
	return nonce, hash[:]
}
 
func (pow *ProofOfWork) Validate() bool {
	var hashInt big.Int
 
	data := pow.prepareData(pow.block.Nonce)
	hash := sha256.Sum256(data)
	hashInt.SetBytes(hash[:])
 
	isValid := hashInt.Cmp(pow.target) == -1
 
	return isValid
}
 
func NewBlock(data string, prevBlockHash []byte) *Block {
	block := &Block{time.Now().Unix(), []byte(data), prevBlockHash, []byte{},0}
	pow := NewProofOfWork(block)
	nonce, hash := pow.Run()
 
	block.Hash = hash[:]
	block.Nonce = nonce
 
	return block
}
 
type Blockchain struct {
	blocks []*Block
}
func (bc *Blockchain) AddBlock(data string) {
	prevBlock := bc.blocks[len(bc.blocks)-1]
	newBlock := NewBlock(data, prevBlock.Hash)
	bc.blocks = append(bc.blocks, newBlock)
}
 
func NewGenesisBlock() *Block {
	return NewBlock("Genesis Block", []byte{})
}
 
func NewBlockchain() *Blockchain {
	return &Blockchain{[]*Block{NewGenesisBlock()}}
}
 
func check(e error) {
	if e != nil {
		panic(e)
	}
}
/**
 * 判断文件是否存在  存在返回 true 不存在返回false
 */
 func checkFileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}

func main() {
	var f *os.File
	var err error
	filename := "./block.txt"
	if checkFileIsExist(filename) { //如果文件存在
		f, err = os.OpenFile(filename, os.O_APPEND, 0666) //打开文件
		fmt.Println("文件存在")
	} else {
		f, err = os.Create(filename) //创建文件
		fmt.Println("文件不存在")
	}
	
	if err != nil {
        fmt.Println(err)
                f.Close()
        return
    }

	bc := NewBlockchain()
	fmt.Println()

	bc.AddBlock("a send 1 yuan to b")
	fmt.Println()
	
	bc.AddBlock("c send 2 yuan to a")
	fmt.Println()
	for _, block := range bc.blocks {
		fmt.Printf("prehash: %x\n", block.PrevBlockHash)
		fmt.Printf("tx: %s\n", block.Data)
		fmt.Printf("Hash: %x\n", block.Hash)
		fmt.Printf("nonce: %x\n", block.Nonce)
		fmt.Println()
		a := string(block.PrevBlockHash)
		a2 :=string(block.Data)
		a3 :=string(block.Hash)
		a4 :=string(block.Nonce)
		d := []string{a,a2,a3,a4}
		for _, v := range d {
			fmt.Fprintln(f, v)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
		fmt.Printf("output: %x\n", d)
	}
	err = f.Close()
}
