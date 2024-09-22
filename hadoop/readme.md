# Hadoop Example

![example](./example.gif)

Rodar o ambiente e o algoritmo Contador de Palavras

```
cd ./src/
./compile.sh
docker-compose build
docker-compose up -d
docker exec -it hadoop_namenode_1 /run.sh
```


## Dowload jdk for compile

```
wget https://download.java.net/openjdk/jdk8u44/ri/openjdk-8u44-linux-x64.tar.gz
```