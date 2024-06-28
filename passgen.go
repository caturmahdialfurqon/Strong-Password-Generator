package main

import (
	"os"
    	"os/exec"
	"crypto/rand"
	"fmt"
	"math/big"
	"strconv"
	"strings"
)

const (
	lowercase       = "abcdefghjkmnpqrstuvwxyz"
	uppercase       = "ABCDEFGHJKMNPQRSTUVWXYZ"
	numbers         = "0123456789"
	symbols         = "!#$%&'()*+,-./:;<=>?@[]^_`{|}~"
	similarChars    = "il1Lo0O"
	sequentialChars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	// just for color
	Cyan    		= "\033[36m"
	Blue    		= "\033[34m" 
	Reset           = "\033[0m" 
)

type PasswordConfig struct {
	Length            int
	IncludeNumbers    bool
	IncludeLowercase  bool
	IncludeUppercase  bool
	IncludeSymbols    bool
	BeginWithLetter   bool
	NoSimilarChars    bool
	NoDuplicateChars  bool
	NoSequentialChars bool
	Quantity          int
}

func generatePassword(config PasswordConfig) (string, error) {
	var charset string
	if config.IncludeLowercase {
		charset += lowercase
	}
	if config.IncludeUppercase {
		charset += uppercase
	}
	if config.IncludeNumbers {
		charset += numbers
	}
	if config.IncludeSymbols {
		charset += symbols
	}

	if config.NoSimilarChars {
		for _, char := range similarChars {
			charset = strings.ReplaceAll(charset, string(char), "")
		}
	}

	if len(charset) == 0 {
		return "", fmt.Errorf("no characters available to generate password")
	}

	password := make([]byte, config.Length)
	usedChars := make(map[byte]bool)

	for i := 0; i < config.Length; i++ {
		if i == 0 && config.BeginWithLetter {
			charset = strings.ReplaceAll(charset, numbers, "")
			charset = strings.ReplaceAll(charset, symbols, "")
		}

		for {
			index, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
			if err != nil {
				return "", err
			}
			char := charset[index.Int64()]

			if config.NoDuplicateChars && usedChars[char] {
				continue
			}

			if config.NoSequentialChars && i > 0 && isSequential(password[i-1], char) {
				continue
			}

			password[i] = char
			usedChars[char] = true
			break
		}
	}

	return string(password), nil
}

func isSequential(prev, curr byte) bool {
	return strings.IndexByte(sequentialChars, prev)+1 == strings.IndexByte(sequentialChars, curr)
}

func createPasswordConfig() (PasswordConfig, error) {
	var length int
	var includeNumbers, includeLowercase, includeUppercase, includeSymbols bool
	var beginWithLetter, noSimilarChars, noDuplicateChars, noSequentialChars bool
	var quantity int

	fmt.Print("Enter desired password length (minimum 6 [weak], [ >16 strong], maximum 50): ")
	for {
		lengthStr, err := readString()
		if err != nil {
			return PasswordConfig{}, err
		}
		length, err = strconv.Atoi(lengthStr)
		if err != nil || length < 6 {
			fmt.Println("Invalid length. Please enter a number greater than or equal to 8.")
			continue
		}
		break
	}

	fmt.Print("Include numbers (y/n)? ")
	includeNumbers, err := readYN()
	if err != nil {
		return PasswordConfig{}, err
	}

	fmt.Print("Include lowercase letters (y/n)? ")
	includeLowercase, err = readYN()
	if err != nil {
		return PasswordConfig{}, err
	}

	fmt.Print("Include uppercase letters (y/n)? ")
	includeUppercase, err = readYN()
	if err != nil {
		return PasswordConfig{}, err
	}

	fmt.Print("Include symbols (y/n)? ")
	includeSymbols, err = readYN()
	if err != nil {
		return PasswordConfig{}, err
	}

	fmt.Print("Begin password with a letter (y/n)? ")
	beginWithLetter, err = readYN()
	if err != nil {
		return PasswordConfig{}, err
	}

	fmt.Print("Exclude similar characters (like 'l' and '1') (y/n)? ")
	noSimilarChars, err = readYN()
	if err != nil {
		return PasswordConfig{}, err
	}

	fmt.Print("Disallow repeated characters in the password (y/n)? ")
	noDuplicateChars, err = readYN()
	if err != nil {
		return PasswordConfig{}, err
	}

	fmt.Print("Prevent consecutive characters from the keyboard layout (y/n)? ")
	noSequentialChars, err = readYN()
	if err != nil {
		return PasswordConfig{}, err
	}

	fmt.Print("Enter the number of passwords to generate: ")
	for {
		quantityStr, err := readString()
		if err != nil {
			return PasswordConfig{}, err
		}
		quantity, err = strconv.Atoi(quantityStr)
		if err != nil || quantity <= 0 {
			fmt.Println("Invalid quantity. Please enter a positive number.")
			continue
		}
		break
	}

	return PasswordConfig{
		Length:           length,
		IncludeNumbers:   includeNumbers,
		IncludeLowercase: includeLowercase,
		IncludeUppercase: includeUppercase,
		IncludeSymbols:   includeSymbols,
		BeginWithLetter:  beginWithLetter,
		NoSimilarChars:   noSimilarChars,
		NoDuplicateChars: noDuplicateChars,
		NoSequentialChars: noSequentialChars,
		Quantity:         quantity,
	}, nil
}

func readString() (string, error) {
	var input string
	_, err := fmt.Scanln(&input)
	return strings.TrimSpace(input), err // Trim leading/trailing whitespace
}

func readYN() (bool, error) {
	var input string
	_, err := fmt.Scanln(&input)
	if err != nil {
		return false, err
	}
	input = strings.ToLower(input) // Convert input to lowercase for case-insensitive check
	switch input {
	case "y", "yes":
		return true, nil
	case "n", "no":
		return false, nil
	default:
		return false, fmt.Errorf("Invalid input. Please enter 'y' or 'n'.")
	}
}


func main() {
	c := exec.Command("clear")
	c.Stdout = os.Stdout
	c.Run()
  config, err := createPasswordConfig()
  if err != nil {
    fmt.Println("Error creating password configuration:", err)
    return
  }
  fmt.Println(Cyan + "Generated passwords:\n" + Reset)
  for i := 0; i < config.Quantity; i++ {
    password, err := generatePassword(config)
    if err != nil {
      fmt.Println("Error generating password:", err)
      continue
    }
    fmt.Println(Blue + password)
  }
}

