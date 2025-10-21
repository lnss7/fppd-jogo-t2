//aqui ta as informações que os outros usuários e o servidor precisa saber!!!

package main

type Jogador struct {
	Nome string
	X    int
	Y    int
}

type Mensagem struct {
	Remetente string
	Texto     string
}

type Movimento struct {
	Nome string
	X    int
	Y    int
}

type EstadoJogo struct {
	Estados []Jogador
}
