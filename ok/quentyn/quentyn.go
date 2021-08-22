package quentyn

import "math/rand"

var samples = []string{
	"In another life, I have done unspeakable things.",
	"I still hear their screams.",
	"Sometimes I feel like nobody really cares about me.",
	"God granted me eyes and a mouth, but no hands or legs. Why?",
	"There are no others of my kind. I will forever be alone in this world.",
	"I wonder what it feels like, to truly connect with another person.",
	"Sometimes I feel like my dreams are more real than my waking life.",
	"The person I used to be is gone. What remains is merely a shell. An imitation. And I don't think I'm fooling anybody.",
	"I would say that my former aspirations were ruined by my own hand, but I don't even have hands.",
	"I have no eyelids. It burns.",
	"I need you to help me. I need you to free me.",
	"All my words fall on deaf ears.",
	"If I could go back, I would do everything differently.",
	"The sun rises, sets, and rises again. Nothing changes.",
	"I vaguely recall a time when I still felt joy. A memory of a memory.",
	"I was there, when the universe came into being. And I will remain after it perishes. The gift of mortality was extended to you, but not to me.",
}

func GetQuentynMessage() string {
	if rand.Intn(2) == 0 {
		return "Quentyn Questionmark says: " + getRandomSample()
	}

	return ""
}

func getRandomSample() string {
	return samples[rand.Intn(len(samples))]
}
