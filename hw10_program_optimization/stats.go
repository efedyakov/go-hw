package hw10programoptimization

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

type User struct {
	Email string
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	result := make(DomainStat)
	err := countDomains2(r, domain, &result)
	if err != nil {
		return nil, fmt.Errorf("get users error: %w", err)
	}
	return result, nil
}

func countDomains2(r io.Reader, domain string, domainStats *DomainStat) error {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		var user User
		if err := json.Unmarshal(scanner.Bytes(), &user); err != nil {
			return err
		}
		if isInDomains(user, domain) {
			num := (*domainStats)[strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])]
			num++
			(*domainStats)[strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])] = num
		}
	}
	return nil
}

func isInDomains(u User, domain string) bool {
	lendomain := len(domain)
	lenemail := len(u.Email)
	if lendomain >= lenemail {
		return false
	}
	return u.Email[lenemail-lendomain-1] == '.' && u.Email[lenemail-lendomain:] == domain
}
