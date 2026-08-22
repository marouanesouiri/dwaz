-- KEYS[1]: bucket key (e.g. "wumpkit:identify_ratelimit:bucket:0")
-- ARGV[1]: interval in milliseconds (e.g. 5000)
local key = KEYS[1]
local interval = tonumber(ARGV[1])

-- get current Redis server time in milliseconds
local redis_time = redis.call('TIME')
local now = tonumber(redis_time[1]) * 1000 + math.floor(tonumber(redis_time[2]) / 1000)

local current_next = redis.call('GET', key)
if not current_next then
    local next_time = now + interval
    redis.call('SET', key, tostring(next_time), 'PX', interval * 2)
    return 0
end

local next_allowed = tonumber(current_next)
if next_allowed <= now then
    local next_time = now + interval
    redis.call('SET', key, tostring(next_time), 'PX', interval * 2)
    return 0
else
    local wait_ms = next_allowed - now
    local next_time = next_allowed + interval
    local ttl = (next_time - now) + interval
    redis.call('SET', key, tostring(next_time), 'PX', ttl)
    return wait_ms
end
