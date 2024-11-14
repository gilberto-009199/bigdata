# Sqoop  Example

Rodar o ambiente e a exportação da tabela livros do postgresql para o hadoop hdfs

```bash
$ docker-compose build --no-cache
$ docker-compose up
$ docker cp src/mapred-site.xml sqoop-nodemanager-1:/opt/hadoop/etc/hadoop/
$ docker exec -it sqoop-nodemanager-1 /bin/sh
$ sqoop list-databases --connect jdbc:postgresql://postgresql:5432/postgres --username userdatabase --password passwrd
$ sqoop list-tables --connect jdbc:postgresql://postgresql:5432/mydatabase  --username userdatabase --password passwrd
$ sqoop import-all-tables --connect jdbc:postgresql://postgresql:5432/mydatabase  --username userdatabase --password passwrd
$ sqoop import --table livros --connect jdbc:postgresql://postgresql:5432/mydatabase  --username userdatabase --password passwrd --target-dir /livros.txt --bindir ./libjars
```
