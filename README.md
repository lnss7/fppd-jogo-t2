# Jogo Multiplayer de Terminal em Go

Este projeto é um jogo multiplayer desenvolvido em Go que roda no terminal usando a biblioteca [termbox-go](https://github.com/nsf/termbox-go). Múltiplos jogadores podem se conectar a um servidor central e compartilhar suas posições e ações em tempo real.

## Arquitetura

O sistema é composto por:
- **Servidor RPC**: Gerencia o estado dos jogadores e processa comandos
- **Cliente**: Interface gráfica e lógica de jogo local
- **Comunicação**: RPC com garantia de execução única (exactly-once)

## Como funciona

- O mapa é carregado de um arquivo `.txt` contendo caracteres que representam diferentes elementos do jogo
- Cada jogador controla um personagem que pode se mover pelo mapa
- Os jogadores veem outros jogadores como símbolos de caveira (☠)
- Toda comunicação é iniciada pelos clientes (pull-based)
- O servidor mantém o estado global e processa comandos com sequence numbers

### Controles

| Tecla | Ação              |
|-------|-------------------|
| W     | Mover para cima   |
| A     | Mover para esquerda |
| S     | Mover para baixo  |
| D     | Mover para direita |
| E     | Interagir         |
| ESC   | Sair do jogo      |

## Instalação e Compilação

1. Instale o Go (versão 1.21 ou superior)
2. Clone este repositório
3. Instale as dependências:

```bash
make deps
```

4. Compile o servidor e cliente:

```bash
make all
```

## Como executar

### Opção 1: Usando Makefile (Recomendado)

1. **Inicie o servidor** (em um terminal):
```bash
make run_servidor
```

2. **Inicie os clientes** (em terminais separados):
```bash
make run_cliente NOME=Joao
make run_cliente NOME=Maria
make run_cliente NOME=Pedro
```

### Opção 2: Executáveis diretos

1. **Inicie o servidor**:
```bash
./servidor
```

2. **Inicie os clientes**:
```bash
./cliente Joao
./cliente Maria
./cliente Pedro
```

### Opção 3: Script de teste automático

```bash
./teste_multiplayer.sh
```

## Estrutura do projeto

### Arquivos principais:
- `servidor.go` — Servidor RPC independente
- `cliente.go` — Cliente multiplayer com interface gráfica
- `tipos.go` — Estruturas de dados para comunicação RPC
- `jogo.go` — Lógica do jogo (mapa, elementos)
- `personagem.go` — Ações do personagem
- `interface.go` — Interface gráfica com termbox

### Arquivos originais (single player):
- `main.go` — Jogo original single player
- `servidorJogo.go` — Implementação RPC básica (não utilizada)

## Características Técnicas

### Garantia de Execução Única (Exactly-Once)
- Cada comando possui um `sequenceNumber` sequencial
- O servidor mantém controle de comandos processados por jogador
- Comandos duplicados são ignorados automaticamente
- Retransmissão de comandos é tratada de forma segura

### Comunicação
- **Protocolo**: RPC (Remote Procedure Call)
- **Padrão**: Pull-based (clientes solicitam atualizações)
- **Frequência**: Atualizações a cada 100ms
- **Porta**: 1234 (configurável)

### Sincronização
- Mutex para acesso thread-safe ao estado dos jogadores
- Thread dedicada para buscar atualizações do servidor
- Estado local sincronizado com servidor

## Exemplo de Uso

1. Abra 4 terminais
2. No terminal 1: `make run_servidor`
3. No terminal 2: `make run_cliente NOME=Jogador1`
4. No terminal 3: `make run_cliente NOME=Jogador2`
5. No terminal 4: `make run_cliente NOME=Jogador3`

Cada jogador verá os outros como símbolos de caveira e poderá se mover pelo mapa compartilhado.

## Troubleshooting

- **Erro de conexão**: Verifique se o servidor está rodando na porta 1234
- **Nome já existe**: Use um nome diferente para o jogador
- **Compilação falha**: Execute `make deps` para instalar dependências
- **Terminal não suporta**: Use um terminal que suporte termbox-go (Linux/macOS recomendados)


