package hw10programoptimization

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

type User struct {
	ID       int    `json:"Id"`
	Name     string `json:"Name"`
	Username string `json:"Username"`
	Email    string `json:"Email"`
	Phone    string `json:"Phone"`
	Password string `json:"Password"`
	Address  string `json:"Address"`
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	u, err := getUsers(r)
	if err != nil {
		return nil, fmt.Errorf("get users error: %w", err)
	}
	return countDomains(u, domain)
}

type users [100_000]User

func getUsers(r io.Reader) (result users, err error) {
	var user User
	scanner := bufio.NewScanner(r)

	for i := 0; scanner.Scan(); i++ {
		if err = user.UnmarshalJSON(scanner.Bytes()); err != nil {
			return
		}
		result[i] = user
	}

	if err = scanner.Err(); err != nil {
		return
	}

	return
}

func countDomains(u users, domain string) (DomainStat, error) {
	result := make(DomainStat)

	for _, user := range u {
		if user.ID == 0 {
			break
		}

		if strings.Contains(user.Email, "."+domain) {
			domainKey := strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])
			result[domainKey]++
		}
	}

	return result, nil
}
