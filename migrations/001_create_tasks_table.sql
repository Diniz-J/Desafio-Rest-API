-- Migration 001: Create tasks table
-- Criação da tabela de tarefas para a TODO API

CREATE TABLE IF NOT EXISTS tasks (
    id VARCHAR(36) PRIMARY KEY COMMENT 'UUID da tarefa',
    title VARCHAR(255) NOT NULL COMMENT 'Título da tarefa (obrigatório)',
    description TEXT COMMENT 'Descrição detalhada (opcional)',
    status ENUM('pending', 'completed') NOT NULL DEFAULT 'pending' COMMENT 'Status atual da tarefa',
    priority ENUM('low', 'medium', 'high') NOT NULL DEFAULT 'medium' COMMENT 'Prioridade da tarefa',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Data de criação',
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'Data da última atualização',
    deleted_at TIMESTAMP NULL DEFAULT NULL COMMENT 'Data de deleção (soft delete - não implementado ainda)',

    INDEX idx_status (status),
    INDEX idx_priority (priority),
    INDEX idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Tabela de tarefas (TODO list)';

-- Criar índice composto para busca por status e prioridade
CREATE INDEX idx_status_priority ON tasks(status, priority);

-- Verificar criação
SHOW TABLES;
DESCRIBE tasks;
