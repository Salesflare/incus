package main

import "sync"

type Storage struct {
    memory      *MemoryStore
    redis       *RedisStore
    StorageType string
    
    userMu      sync.Mutex
    pageMu      sync.Mutex
}

func initStore(Config *Configuration) *Storage {
    store_type := "memory"
    var redisStore *RedisStore
    
    redis_enabled := Config.Get("redis_enabled");
    if(redis_enabled == "true") {
        redis_host := Config.Get("redis_host")
        redis_port := uint(Config.GetInt("redis_port"))
        
        redisStore = newRedisStore(redis_host, redis_port)
        store_type = "redis"
    }
    
    var Store = Storage{
        &MemoryStore{make(map[string]map[string] *Socket), make(map[string] map[string] *Socket), 0},
        redisStore,
        store_type,
        
        sync.Mutex{},
        sync.Mutex{},
    }
    
    return &Store
}

func (this *Storage) Save(sock *Socket) (error) {
    this.userMu.Lock()
    
    this.memory.Save(sock)
    
    if this.StorageType == "redis" {
        if err := this.redis.Save(sock); err != nil {
            return err
        }
    }
    
    this.userMu.Unlock()
    return nil
}

func (this *Storage) Remove(sock *Socket) (error) {
    this.userMu.Lock()
    
    this.memory.Remove(sock)
    
    if this.StorageType == "redis" {
        if err := this.redis.Remove(sock); err != nil {
            return err
        }
    }
    
    this.userMu.Unlock()
    return nil
}

func (this *Storage) Client(UID string) (map[string]*Socket, error) {
    return this.memory.Client(UID)
}

func (this *Storage) Clients() (map[string]map[string]*Socket) {
    return this.memory.Clients()
}

func (this *Storage) ClientList() ([]string, error) {
    if this.StorageType == "redis" {
        return this.redis.Clients()
    }
    
    return nil, nil
}

func (this *Storage) Count() (int64, error) {
    if this.StorageType == "redis" {
        return this.redis.Count()
    }
    
    return this.memory.Count()
}

func (this *Storage) SetPage(sock *Socket) error {
    this.pageMu.Lock()
    
    this.memory.SetPage(sock)
    
    if this.StorageType == "redis" {
        if err := this.redis.SetPage(sock); err != nil {
            return err
        }
    }
    
    this.pageMu.Unlock()
    return nil
}

func (this *Storage) UnsetPage(sock *Socket) error {
    this.pageMu.Lock()
    this.memory.UnsetPage(sock)
    
    if this.StorageType == "redis" {
        if err := this.redis.UnsetPage(sock); err != nil {
            return err
        }
    }
    
    this.pageMu.Unlock()
    return nil
}

func (this *Storage) getPage(page string) map[string]*Socket {
    return this.memory.getPage(page)
}
