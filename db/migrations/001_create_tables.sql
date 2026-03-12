CREATE TABLE IF NOT EXISTS usuarios (
    id SERIAL PRIMARY KEY,
    nome VARCHAR(100) NOT NULL,
    email VARCHAR(150) UNIQUE NOT NULL,
    senha_hash TEXT NOT NULL,
    criado_em TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS reembolsos (
    id SERIAL PRIMARY KEY,
    usuario_id INT NOT NULL REFERENCES usuarios(id),
    descricao TEXT NOT NULL,
    valor NUMERIC(10, 2) NOT NULL,
    categoria VARCHAR(50) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'PENDENTE',
    comprovante TEXT,
    criado_em TIMESTAMP DEFAULT NOW(),
    atualizado_em TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_reembolsos_usuario_id ON reembolsos(usuario_id);
CREATE INDEX idx_reembolsos_status ON reembolsos(status);
