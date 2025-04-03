box.cfg{
    listen = 3301,
    replication = nil, -- Отключение репликации
    read_only = false, -- Разрешение записи
    wal_mode = "write" -- Включение упреждающего журналирования
}

-- Создаем пользователя, если его нет
if not box.schema.user.exists('user') then
    box.schema.user.create('user', { password = 'secret' })
    
    -- Даем права на все операции
    box.schema.user.grant('user', 'read,write,execute', 'universe')
    
    print("User 'user' created!")
else
    print("User' already exists")
end

-- Создание пространтсва 'polls'
if not box.space.polls then
    box.schema.space.create('polls', {
        format = {
            {name = 'id', type = 'string'},
            {name = 'question', type = 'string'},
            {name = 'options', type = 'map'},
            {name = 'voters', type = 'map'},
            {name = 'creator', type = 'string'},
            {name = 'closed', type = 'boolean'}
        },
        if_not_exists = true
    })
    print("Space 'polls' created")
else
    print("Space 'polls' already exists")
end

-- Создание индексов для пространства 'polls'
box.space.polls:create_index('primary', { 
    parts = { {field = 'id', type = 'string'} }, 
    type = 'hash',
    if_not_exists = true 
})

-- Создание пространтсва 'cmd_tokens' (для валидации токенов в "memory" режиме)
if not box.space.cmd_tokens then
    box.schema.space.create('cmd_tokens', {
        format = {
            {name = 'cmd_path', type = 'string'},
            {name = 'token', type = 'string'},
        },
        if_not_exists = true
    })
    print("Space 'cmd_tokens' created")
else
    print("Space 'cmd_tokens' already exists")
end

-- Создание индексов для пространства 'cmd_tokens'
box.space.cmd_tokens:create_index('primary', { 
    parts = { {field = 'cmd_path', type = 'string'} }, 
    type = 'hash',
    if_not_exists = true 
})
