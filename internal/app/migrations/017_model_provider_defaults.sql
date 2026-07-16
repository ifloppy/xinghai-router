insert into model_providers(name, slug, prefixes, priority) values
    ('OpenAI', 'openai', '["gpt-", "o1", "o3", "o4"]', 10),
    ('Anthropic', 'anthropic', '["claude"]', 20),
    ('Google', 'google', '["gemini"]', 30),
    ('DeepSeek', 'deepseek', '["deepseek"]', 40),
    ('Alibaba', 'alibaba', '["qwen", "qwq"]', 50),
    ('Zhipu', 'zhipu', '["glm"]', 60),
    ('Mistral', 'mistral', '["mistral", "codestral"]', 70),
    ('xAI', 'xai', '["grok"]', 80),
    ('Meta', 'meta', '["llama"]', 90),
    ('Moonshot', 'moonshot', '["moonshot", "kimi"]', 100)
on conflict (name) do update set
    slug = excluded.slug,
    prefixes = excluded.prefixes,
    priority = excluded.priority,
    updated_at = now();
