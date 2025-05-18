package utils

import (
	"fmt"
	"hash/fnv"
	"math/rand"
)

var adjectives = []string{
	"Adventurous", "Agile", "Amazing", "Ancient", "Angelic", "Artistic", "Bold",
	"Brave", "Bright", "Calm", "Clever", "Colorful", "Cool", "Creative",
	"Curious", "Daring", "Dazzling", "Eager", "Elegant", "Energetic", "Famous",
	"Fearless", "Feisty", "Fierce", "Friendly", "Fun", "Funny", "Gentle",
	"Glorious", "Grumpy", "Happy", "Heroic", "Hopeful", "Humble", "Hyper",
	"Joyful", "Kind", "Legendary", "Lively", "Lucky", "Luminous", "Majestic",
	"Mighty", "Mischievous", "Mysterious", "Nimble", "Peaceful", "Playful",
	"Proud", "Quick", "Quiet", "Quirky", "Radiant", "Relaxed", "Savage",
	"Shiny", "Silent", "Silly", "Smart", "Sneaky", "Swift", "Tricky",
	"Valiant", "Wild", "Wise", "Witty", "Zany",
}

var animals = []string{
	"Alpaca", "Antelope", "Armadillo", "Aardvark", "Axolotl", "Badger", "Bat",
	"Bear", "Beaver", "Bison", "Buffalo", "Camel", "Capybara", "Caracal",
	"Cat", "Chameleon", "Cheetah", "Chimpanzee", "Chinchilla", "Cobra", "Cougar",
	"Coyote", "Crab", "Crocodile", "Crow", "Deer", "Dingo", "Dolphin",
	"Donkey", "Duck", "Eagle", "Elephant", "Ferret", "Flamingo", "Fox",
	"Frog", "Gazelle", "Giraffe", "Goat", "Gorilla", "Hedgehog", "Hippo",
	"Hyena", "Iguana", "Jackal", "Jaguar", "Kangaroo", "Koala", "Lemur",
	"Leopard", "Lion", "Llama", "Lynx", "Manatee", "Mole", "Monkey",
	"Moose", "Narwhal", "Newt", "Octopus", "Opossum", "Orca", "Otter",
	"Owl", "Panda", "Panther", "Peacock", "Pelican", "Penguin", "Platypus",
	"Porcupine", "Possum", "Quokka", "Rabbit", "Raccoon", "Rat", "Reindeer",
	"Rhino", "Seal", "Serval", "Shark", "Sheep", "Skunk", "Sloth", "Snail",
	"Snake", "Swan", "Tiger", "Toad", "Tortoise", "Turtle", "Walrus",
	"Weasel", "Whale", "Wolf", "Wombat", "Yak", "Zebra",
}

// hashStringToInt64 hashes a string into a deterministic int64 using FNV
func hashStringToInt64(s string) int64 {
	h := fnv.New64a()
	h.Write([]byte(s))
	return int64(h.Sum64())
}

// GenerateName returns a username based on a given seed input (UUID, IP, etc.)
func GenerateName(seedInput string) string {
	seed := hashStringToInt64(seedInput)
	r := rand.New(rand.NewSource(seed))

	adj := adjectives[r.Intn(len(adjectives))]
	animal := animals[r.Intn(len(animals))]

	return fmt.Sprintf("%s%s", adj, animal)
}
