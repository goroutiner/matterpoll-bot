box.cfg{
    listen = 3301,
    replication = nil, -- Отключение репликации
    read_only = false  -- Разрешение записи
}

-- Создаем пользователя, если его нет
if not box.schema.user.exists('user') then
    box.schema.user.create('user', { password = 'secret' })
    
    -- Даем права на все операции
    box.schema.user.grant('user', 'read,write,execute', 'universe')
    
    print("Пользователь 'user' создан!")
else
    print("Пользователь 'user' уже существует.")
end

-- Create a space --
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

-- Create indexes --
box.space.polls:create_index('poll_id_index', { parts = { {field = 'id', type = 'string'} }, if_not_exists = true })
print("\"poll_id_index\" created for space 'polls'")