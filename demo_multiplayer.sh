#!/bin/bash

# Script de demonstração do jogo multiplayer
echo "=== DEMONSTRAÇÃO JOGO MULTIPLAYER ==="
echo ""

# Verifica se o servidor está rodando
if ! pgrep -f "./servidor" > /dev/null; then
    echo "Iniciando servidor..."
    ./servidor &
    SERVIDOR_PID=$!
    sleep 2
    echo "Servidor iniciado (PID: $SERVIDOR_PID)"
else
    echo "Servidor já está rodando"
fi

echo ""
echo "=== TESTE DE COMUNICAÇÃO RPC ==="
echo "Testando conexão e funcionalidades básicas..."
./teste_rpc JogadorTeste
echo ""

echo "=== DEMONSTRAÇÃO MULTIPLAYER ==="
echo "Agora você pode testar com múltiplos clientes:"
echo ""
echo "Terminal 1: ./cliente_texto Joao"
echo "Terminal 2: ./cliente_texto Maria" 
echo "Terminal 3: ./cliente_texto Pedro"
echo ""
echo "Comandos em cada cliente:"
echo "  w/a/s/d - mover"
echo "  e - interagir"
echo "  q - sair"
echo ""
echo "Pressione Enter para continuar ou Ctrl+C para sair..."
read

echo "Iniciando demonstração automática..."
echo ""

# Inicia alguns clientes de demonstração
echo "Iniciando cliente 'Joao'..."
./cliente_texto Joao &
JOAO_PID=$!

sleep 2

echo "Iniciando cliente 'Maria'..."
./cliente_texto Maria &
MARIA_PID=$!

sleep 2

echo "Iniciando cliente 'Pedro'..."
./cliente_texto Pedro &
PEDRO_PID=$!

echo ""
echo "Clientes iniciados! Verifique os terminais para ver a interação."
echo "Pressione Enter para finalizar a demonstração..."
read

# Finaliza os clientes
echo "Finalizando demonstração..."
kill $JOAO_PID $MARIA_PID $PEDRO_PID 2>/dev/null

echo "Demonstração concluída!"
