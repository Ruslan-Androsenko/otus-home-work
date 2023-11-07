package hw10programoptimization

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

type User struct {
	Email string
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
		if i >= len(result) {
			break
		}

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
		if strings.HasSuffix(user.Email, domain) && strings.Contains(user.Email, "@") {
			domainKey := strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])
			result[domainKey]++
		}
	}

	return result, nil
}
