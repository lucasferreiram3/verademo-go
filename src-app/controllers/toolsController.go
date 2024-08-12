package controllers

import (
	"net/http"
	"os/exec"
	"verademo-go/src-app/shared/view"

	"golang.org/x/exp/rand"
)

type Outputs struct {
	Host    string
	Ping    string
	Fortune string
}

func ShowTools(w http.ResponseWriter, r *http.Request) {
	view.Render(w, "tools.html", nil)
}

func ProcessTools(w http.ResponseWriter, r *http.Request) {
	host := r.FormValue("host")
	fortuneFile := r.FormValue("fortunefile")

	var outputs Outputs

	if host == "" {
		outputs.Ping = ""
	} else {
		outputs.Ping = ping(host)
		outputs.Host = host
	}

	if fortuneFile == "" {
		outputs.Fortune = ""
	} else {
		outputs.Fortune = fortune(fortuneFile)
	}

	view.Render(w, "tools.html", outputs)
}

func ping(host string) string {
	ping, err := exec.Command("ping", "-c", "1", host).Output()
	var pingResult string
	if err != nil {
		pingResult = "Error pinging host " + host + ". Please check the hostname."
	} else {
		pingResult = string(ping)
	}
	return pingResult
}

func fortune(fortuneFile string) string {
	fortune, err := exec.Command("fortune", fortuneFile).Output()
	var fortuneResult string
	if err != nil {
		if fortuneFile == "fortune" {
			fortuneResult = fortunes[rand.Intn(len(fortunes))]
		} else if fortuneFile == "riddle" {
			fortuneResult = riddles[rand.Intn(len(riddles))]
		} else {
			fortuneResult = "No valid fortune found."
		}
	} else {
		fortuneResult = string(fortune)
	}
	return fortuneResult
}

var fortunes = [...]string{
	"You will have a great day today.",
	"A new opportunity will come your way.",
	"Someone special is thinking of you.",
	"Expect good news in the mail.",
	"You will meet someone who will change your life.",
	"A financial opportunity will present itself.",
	"You will find what you have been seeking.",
	"A fresh start will put you on your way.",
	"An unexpected event will bring you joy.",
	"Believe in yourself and others will too.",
	"Your hard work will soon pay off.",
	"You will travel to a new place soon.",
	"Your kindness will lead you to success.",
	"Something you lost will soon turn up.",
	"You will soon receive a pleasant surprise.",
	"A new friendship will bring joy to your life.",
	"You are capable of achieving greatness.",
	"Good things are coming your way.",
	"You will overcome a great challenge.",
	"Someone will offer you helpful advice.",
	"You will accomplish your goals.",
	"A creative idea will bring you success.",
	"You will soon be surrounded by good friends.",
	"Your positive attitude will lead you to success.",
	"You will receive a compliment from someone important.",
	"Your efforts will soon be rewarded.",
	"You will make a positive difference in someone's life.",
	"A small act of kindness will go a long way.",
	"You will discover a new talent within yourself.",
	"Your future is bright.",
	"A surprise gift is coming your way.",
	"You will have a reason to celebrate soon.",
	"Your dreams will come true.",
	"Someone will express their gratitude to you.",
	"You will find clarity in a confusing situation.",
	"A new perspective will lead to a great opportunity.",
	"You will make a valuable discovery.",
	"Your good deeds will be repaid.",
	"You will experience something wonderful today.",
	"You will achieve your goals sooner than you think.",
	"Someone will help you achieve your goals.",
	"You will be surrounded by positivity.",
	"Your hard work will be recognized.",
	"You will find peace in a difficult situation.",
	"A long-awaited event will bring joy to your life.",
	"You will meet someone who will inspire you.",
	"You will receive good news from afar.",
	"Your positive outlook will lead to great things.",
	"You will have a happy encounter.",
	"You will achieve something you thought was impossible.",
	"Your kindness will be returned in unexpected ways.",
	"You will find a solution to a problem.",
	"You will experience a stroke of luck.",
	"Your generosity will be rewarded.",
	"You will soon be surrounded by good fortune.",
	"You will make a new friend.",
	"You will receive praise for your efforts.",
	"You will be successful in your endeavors.",
	"You will find joy in unexpected places.",
	"You will have a reason to smile today.",
	"A happy event is in your near future.",
	"You will discover something new about yourself.",
	"You will achieve something you have been working towards.",
	"You will have a positive impact on someone's life.",
	"You will experience a moment of clarity.",
	"A new opportunity will bring you success.",
	"You will find happiness in small things.",
	"You will receive unexpected help.",
	"You will accomplish more than you think.",
	"You will have a pleasant surprise today.",
	"You will find what you have been missing.",
	"You will soon achieve a long-term goal.",
	"Your efforts will lead to success.",
	"You will be recognized for your hard work.",
	"You will find peace and happiness.",
	"You will make a significant achievement.",
	"You will experience a moment of joy.",
	"You will find success in your endeavors.",
	"You will discover a new passion.",
	"You will have a reason to be proud.",
	"You will make a valuable connection.",
	"You will receive a pleasant surprise.",
	"You will be appreciated for your efforts.",
	"You will have a positive experience today.",
	"You will achieve something great.",
	"You will be surrounded by happiness.",
	"You will find joy in helping others.",
	"You will achieve a personal milestone.",
	"You will receive a compliment that makes your day.",
	"You will have a reason to celebrate.",
	"You will accomplish something significant.",
	"You will find what you are looking for.",
	"You will make a positive impact on someone.",
	"You will discover something valuable.",
	"You will receive good news today.",
	"You will have a memorable experience.",
	"You will be recognized for your kindness.",
	"You will find joy in unexpected moments.",
	"You will achieve your dreams.",
	"You will have a reason to be happy.",
	"You will find success in your journey.",
	"You will receive a pleasant message.",
}

var riddles = [...]string{
	"Q: Why haven't you graduated yet?\nA: Well, Dad, I could have finished years ago, but I wanted my dissertation to rhyme.",
	"Q: Why did the cow cross the road?\nA: To get to the other side.",
	"Q: What has keys but can't open locks?\nA: A piano.",
	"Q: What has a heart that doesn't beat?\nA: An artichoke.",
	"Q: What can travel around the world while staying in a corner?\nA: A stamp.",
	"Q: What comes once in a minute, twice in a moment, but never in a thousand years?\nA: The letter 'M'.",
	"Q: I am not alive, but I can grow; I don't have lungs, but I need air; I don't have a mouth, but water kills me. What am I?\nA: Fire.",
	"Q: I have branches, but no fruit, trunk, or leaves. What am I?\nA: A bank.",
	"Q: What has to be broken before you can use it?\nA: An egg.",
	"Q: I'm tall when I'm young, and I'm short when I'm old. What am I?\nA: A candle.",
	"Q: What month of the year has 28 days?\nA: All of them.",
	"Q: What is full of holes but still holds water?\nA: A sponge.",
	"Q: What question can you never answer yes to?\nA: Are you asleep yet?",
	"Q: What is always in front of you but can't be seen?\nA: The future.",
	"Q: There's a one-story house in which everything is yellow. Yellow walls, yellow doors, yellow furniture. What color are the stairs?\nA: There aren't any—it's a one-story house.",
	"Q: What can you break, even if you never pick it up or touch it?\nA: A promise.",
	"Q: What goes up but never comes down?\nA: Your age.",
	"Q: A man who was outside in the rain without an umbrella or hat didn't get a single hair on his head wet. Why?\nA: He was bald.",
	"Q: What gets wet while drying?\nA: A towel.",
	"Q: What can you keep after giving to someone?\nA: Your word.",
	"Q: I shave every day, but my beard stays the same. What am I?\nA: A barber.",
	"Q: You see me once in June, twice in November, and not at all in May. What am I?\nA: The letter 'E'.",
	"Q: I have lakes with no water, mountains with no stone, and cities with no buildings. What am I?\nA: A map.",
	"Q: What is seen in the middle of March and April that can't be seen at the beginning or end of either month?\nA: The letter 'R'.",
	"Q: You see a boat filled with people. It has not sunk, but when you look again you don't see a single person on the boat. Why?\nA: All the people were married.",
	"Q: What runs all around a backyard, yet never moves?\nA: A fence.",
	"Q: What can fill a room but takes up no space?\nA: Light.",
	"Q: If you drop me, I'm sure to crack, but give me a smile and I'll always smile back. What am I?\nA: A mirror.",
	"Q: I have keys but no locks. I have a space but no room. You can enter, but you can't go outside. What am I?\nA: A keyboard.",
	"Q: What has many teeth but can't bite?\nA: A comb.",
	"Q: Where does today come before yesterday?\nA: In a dictionary.",
	"Q: What kind of band never plays music?\nA: A rubber band.",
	"Q: What has lots of eyes, but can't see?\nA: A potato.",
	"Q: What has one eye, but can't see?\nA: A needle.",
	"Q: What has many needles, but doesn't sew?\nA: A Christmas tree.",
	"Q: What is cut on a table, but is never eaten?\nA: A deck of cards.",
	"Q: What has a head, a tail, is brown, and has no legs?\nA: A penny.",
	"Q: What is black when it's clean and white when it's dirty?\nA: A chalkboard.",
	"Q: What gets bigger the more you take away?\nA: A hole.",
	"Q: I'm light as a feather, yet the strongest man can't hold me for much longer. What am I?\nA: Breath.",
	"Q: I'm found in socks, scarves and mittens; and often in the paws of playful kittens. What am I?\nA: Yarn.",
	"Q: Where does one wall meet the other wall?\nA: On the corner.",
	"Q: What building has the most stories?\nA: The library.",
	"Q: What tastes better than it smells?\nA: Your tongue.",
	"Q: What has one head, one foot, and four legs?\nA: A bed.",
	"Q: What has many keys but can't open a single lock?\nA: A piano.",
	"Q: What begins with T, ends with T, and has T in it?\nA: A teapot.",
	"Q: What has words, but never speaks?\nA: A book.",
	"Q: What has four wheels and flies?\nA: A garbage truck.",
}
