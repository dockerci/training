package main

import (
	"fmt"
	"encoding/hex"
	"net/http"
	"log"
	"encoding/binary"
	"crypto/rand"
	"os"
)

func query(ciphertext string) int{
	HTTP_request := "http://crypto-class.appspot.com/po?er=" + ciphertext
/
	resp, err := http.Get(HTTP_request)
	if err != nil {	
		fmt.Printf("Can not retrieve request")
	}

	defer resp.Body.Close()
	// body, err := ioutil.ReadAll(resp.Body)
	statusRes := resp.StatusCode
	return statusRes	
}

func decrypOracle(block [16]byte,blockId int) [16]byte{
	const b int = 16
	var (
		tmp byte
		count byte
		err_num int
		result [b]byte 
		requestBytes [2*b]byte 
		j,i int
	)
	for iter := b ; iter < 2*b ; iter++ {
		requestBytes[iter] = block[iter - b]
	}

	for j = 1 ; j < b ; j++ {
		binary.Read(rand.Reader,binary.LittleEndian, &requestBytes[j])	
	}
	fmt.Printf("Guess byte : 15 at block %d\n",blockId)
	count = 0
	requestBytes[b - 1] = requestBytes[b-1] ^ count
	err_num = query(hex.EncodeToString(requestBytes[:]))
	
	for err_num != 404 {
		//fmt.Printf("%d : %d\n",err_num,count)
		count++
		if count > 255 {
			fmt.Printf("We are fail !!!")
			os.Exit(1)
		}
		requestBytes[b-1] = requestBytes[b-1] ^ (count - 1) ^ count
		err_num = query(hex.EncodeToString(requestBytes[:]))

	}

	result[b - 1] = requestBytes[b-1] ^ 1 
	fmt.Printf("Done\n")
	for j = b - 2 ; j >= 0 ; j-- {
		fmt.Printf("Guess byte : %d at block %d\n",j,blockId)	
		tmp = byte(b - j)
		count = 0
		
		for i = j + 1 ; i < b  ; i++ {
			requestBytes[i] = result[i] ^ tmp
		}
	
		requestBytes[j] = requestBytes[j] ^ count
		err_num = query(hex.EncodeToString(requestBytes[:]))
		for err_num != 404 {
			//fmt.Printf("%d : %d\n",err_num,count)
			count++
			if count > 255 {
				fmt.Printf("We are fail !!!")
				os.Exit(1)
			}
			requestBytes[j] = requestBytes[j] ^ (count-1) ^ count
			err_num = query(hex.EncodeToString(requestBytes[:]))
		}
		fmt.Printf("Done\n")
		result[j] = requestBytes[j] ^ tmp
	}

	return result
}

func StringToArrayOfBlock(cipher string) ([][16]byte ,int){
	aBytes ,error := hex.DecodeString(cipher)
	if error != nil {
		log.Fatal(error)
	}
	numBlock := len(aBytes)/16
	slice := make([][16]byte , numBlock) 
	
	for i := 0 ; i < numBlock ; i++ {
		var newBlock [16]byte
		for j := 0 ; j < 16 ; j++ {
			newBlock[j] = aBytes[i*16 + j]
		}
		slice[i] = newBlock

	}
	return slice , numBlock
}


func decrypt(cipher string) string {
	//cBytes ,error:= hex.DecodeString(cipher)
	slice , numBlocks := StringToArrayOfBlock(cipher)
	resultSlice := make([][16]byte , numBlocks)
	decryptBytes := make([]byte,(numBlocks - 1)*16)
	for i := numBlocks - 1 ; i > 0 ; i-- {
		fmt.Printf("Block %d \n",i)
		resultSlice[i] = decrypOracle(slice[i],i) 
	}	
	//resultSlice[0] = slice[0]

	for i := numBlocks - 1; i > 0 ; i-- {
		for j := 0 ; j < 16 ; j++ {
			decryptBytes[(i-1)*16 + j] = resultSlice[i][j] ^ slice[i - 1][j] 
		}
	}

	plainText := string(decryptBytes)	
	
	return plainText 
}




func main() {
	
	a := "f20bdba6ff29eed7b046d1df9fb7000058b1ffb4210a580f748b4ac714c001bd4a61044426fb515dad3f21f18aa577c0bdf302936266926ff37dbf7035d5eeb4"	
	b := decrypt(a)
	fmt.Printf("%s\n",b)
	fmt.Printf("Credit : p77u4n - 2015\n")
	
}
