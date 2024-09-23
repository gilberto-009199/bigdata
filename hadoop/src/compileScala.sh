#!/bin/sh

# Caminho onde o script está localizado (diretório atual do script)
PWDCOMPILE="$(cd "$(dirname "$0")" && pwd)"

# Compilar arquivos Scala
../scala-2.13.0/bin/scalac -cp "$PWDCOMPILE/hadoop-client-3.4.0.jar:$PWDCOMPILE/hadoop-common-3.4.0.jar:$PWDCOMPILE/hadoop-hdfs-3.4.0.jar:$PWDCOMPILE/hadoop-mapreduce-client-core-3.4.0.jar" *.scala

# Verificar se a compilação do Scala foi bem-sucedida
if [ $? -eq 0 ]; then

    # Verificar se a compilação do Java foi bem-sucedida
    if [ $? -eq 0 ]; then
        echo "Compilação bem-sucedida!"
        echo "Empacotando....."
        
        # Empacotar em um JAR
        ../java-se-8u44-ri/bin/jar cf $PWDCOMPILE/contadorDePalavras.jar *.class

        # Verificar se o empacotamento foi bem-sucedido
        if [ $? -eq 0 ]; then
            echo "Empacotamento bem-sucedido: contadorDePalavras.jar"
        else
            echo "Erro no empacotamento!"
        fi
    else
        echo "Erro na compilação do Scala!"
    fi
else
    echo "Erro na compilação do Scala!"
fi
