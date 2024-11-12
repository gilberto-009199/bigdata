-- Criação da tabela de autores
CREATE TABLE autores (
    autor_id SERIAL PRIMARY KEY,
    nome VARCHAR(100) NOT NULL,
    nacionalidade VARCHAR(50),
    data_nascimento DATE
);

-- Criação da tabela de livros
CREATE TABLE livros (
    livro_id SERIAL PRIMARY KEY,
    titulo VARCHAR(150) NOT NULL,
    genero VARCHAR(50),
    ano_publicacao INT,
    resumo TEXT
);

-- Tabela associativa para a relação muitos-para-muitos entre livros e autores
CREATE TABLE livros_autores (
    livro_id INT REFERENCES livros(livro_id) ON DELETE CASCADE,
    autor_id INT REFERENCES autores(autor_id) ON DELETE CASCADE,
    PRIMARY KEY (livro_id, autor_id)
);

-- Inserindo autores
INSERT INTO autores (nome, nacionalidade, data_nascimento)
VALUES
('Gabriel Garcia Marquez', 'Colombiana', '1927-03-06'),
('J.K. Rowling', 'Britânica', '1965-07-31'),
('George R.R. Martin', 'Americana', '1948-09-20');

-- Inserindo livros
INSERT INTO livros (titulo, genero, ano_publicacao, resumo)
VALUES
('Cem Anos de Solidão', 'Romance', 1967, 'Um épico da família Buendía na aldeia fictícia de Macondo.'),
('Harry Potter e a Pedra Filosofal', 'Fantasia', 1997, 'O primeiro livro da série Harry Potter.'),
('A Guerra dos Tronos', 'Fantasia', 1996, 'Primeiro livro da série de fantasia épica As Crônicas de Gelo e Fogo.');

-- Relacionando livros e autores
INSERT INTO livros_autores (livro_id, autor_id)
VALUES
(1, 1),  -- Cem Anos de Solidão por Gabriel Garcia Marquez
(2, 2),  -- Harry Potter e a Pedra Filosofal por J.K. Rowling
(3, 3);  -- A Guerra dos Tronos por George R.R. Martin
