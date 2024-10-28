CREATE TABLE job_offers (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT NOT NULL,
    author TEXT NOT NULL,
    description TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    views_count INTEGER DEFAULT 0
);

CREATE TABLE views (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    job_offer_id INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO job_offers (id, title, author, description) VALUES
(1, 'Senior Software Engineer', 'TechCorp Inc.', 'Join our team to build cutting-edge web applications using the latest technologies.'),
(2, 'UX Designer', 'DesignMasters LLC', 'Create intuitive and beautiful user interfaces for our clients'' products.'),
(3, 'Senior Unicorn Wrangler', 'MagicTech Inc.', 'Seeking an experienced unicorn handler to manage our growing herd of mythical creatures. Must be proficient in rainbow magic and horn polishing.'),
(4, 'Chief Meme Officer', 'Viral Dreams LLC', 'Lead our meme creation team to produce the dankest content on the internet. Knowledge of current trends and cat videos is a must.');

create virtual table job_offers_embeddings using vec0(
    embedding float[384]
);
