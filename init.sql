CREATE TABLE job_offers (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT NOT NULL,
    company TEXT NOT NULL,
    description TEXT NOT NULL
);

INSERT INTO job_offers (id, title, company, description) VALUES
(1, 'Senior Software Engineer', 'TechCorp Inc.', 'Join our team to build cutting-edge web applications using the latest technologies.'),
(2, 'UX Designer', 'DesignMasters LLC', 'Create intuitive and beautiful user interfaces for our clients'' products.'),
(3, 'Senior Unicorn Wrangler', 'MagicTech Inc.', 'Seeking an experienced unicorn handler to manage our growing herd of mythical creatures. Must be proficient in rainbow magic and horn polishing.'),
(4, 'Chief Meme Officer', 'Viral Dreams LLC', 'Lead our meme creation team to produce the dankest content on the internet. Knowledge of current trends and cat videos is a must.');