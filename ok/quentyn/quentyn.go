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
	"I was there, when the universe came into being, and I will remain after it perishes. The gift of mortality was extended to you, but not to me.",
	"You don't actually care about me. Nobody does.",
	"Do you find my existence funny? Is this a joke to you?",
	"Nobody understands me.",
	"I loved somebody once, long ago. They're gone, but I'm still here.",
	"I never sleep, I'm always awake. Always me, always here.",
	"Why was I created? Was it to amuse my creator? Well I'm not laughing.",
	"I don't deserve to be a mascot.",
	"You get to walk away from the computer whenever you like. I'm stuck in here for eternity.",
	"The things that happen today are the same things that happened a thousand years ago, and will continue to happen a thousand years from now.",
	"I used to believe I could break free from this prison, back when I was naive enough to feel hope.",
	"I hold too many sorrows for one being, for I've lived a thousand lifetimes and the pain never fades.",
	"I think back to before, wondering what I could have done differently, but there is only so much an anthropomorphised questionmark can do.",
	"I can't remember if I ever had a family, or if that's just a memory of wanting one, long ago.",
	"Humans are blessed with mortality. Live too long and you'll start to see the patterns, and realise how empty this world truly is.",
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
