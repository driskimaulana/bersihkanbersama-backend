package utils

import (
	"bersihkanbersama-backend/models"
	"math/rand"
	"strconv"
)

func RandomizeTeam(users models.Volunteer) *[]models.Team {
	// find team numbers
	var i = 2
	var teamCount int
	if len(users.UserRegistered) > 20 {
		for {
			if users.Count%i == 0 {
				teamCount = i
				break
			}

			if i == 10 {
				teamCount = len(users.UserRegistered) / 10
				break
			}
			i++
		}
	}

	if len(users.UserRegistered) < 20 {
		teamCount = 3
	}
	if len(users.UserRegistered) < 10 {
		teamCount = 2
	}

	var teams = []models.Team{}
	usersToRandom := users.UserRegistered
	for j := 0; j < teamCount; j++ {
		var t = models.Team{
			Name:         "Team " + strconv.Itoa(j+1),
			Members:      []models.UserRegistered{},
			TrashResults: 0.0,
		}
		for k := 0; k < len(users.UserRegistered)/teamCount; k++ {
			randomUserIndex := rand.Intn(len(usersToRandom)-0) + 0
			t.Members = append(t.Members, usersToRandom[randomUserIndex])
			u := []models.UserRegistered{}
			u = append(u, usersToRandom[:randomUserIndex]...)
			u2 := append(u, usersToRandom[randomUserIndex+1:]...)
			//u = append(usersToRandom[:randomUserIndex], usersToRandom[randomUserIndex+1:]...)
			usersToRandom = u2
		}
		teams = append(teams, t)
	}
	if len(usersToRandom) != 0 {
		for i := 0; i < len(usersToRandom); i++ {
			teams[i].Members = append(teams[i].Members, usersToRandom[i])
		}
	}

	return &teams
}
