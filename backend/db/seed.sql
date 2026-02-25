-- Seed Users (passwords are bcrypt hash of 'password123')
INSERT INTO users (username, email, phone, password, skills, role) VALUES
('TechCorp Kenya', 'hr@techcorp.co.ke', '+254700100100', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', NULL, 'employer'),
('GreenFarms Ltd', 'jobs@greenfarms.co.ke', '+254700200200', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', NULL, 'employer');

INSERT INTO users (username, email, phone, password, skills, role) VALUES
('Amina Wanjiku', 'amina@example.com', '+254711000001', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'javascript,react,nodejs', 'jobseeker'),
('Brian Odhiambo', 'brian@example.com', '+254711000002', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'python,data-analysis,excel', 'jobseeker'),
('Cynthia Mwende', 'cynthia@example.com', '+254711000003', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'graphic-design,photoshop,figma', 'jobseeker');

-- Seed Jobs
INSERT INTO jobs (employer_id, title, description, skills, salary, location, status) VALUES
(1, 'Junior Frontend Developer', 'Build responsive web interfaces using React and modern CSS. Great opportunity for recent graduates.', 'javascript,react,css', 'KES 50,000 - 80,000', 'Nairobi', 'open'),
(1, 'Data Entry Clerk', 'Accurate data entry and record keeping for our operations team. No experience required.', 'typing,excel,attention-to-detail', 'KES 25,000 - 35,000', 'Nairobi', 'open'),
(1, 'IT Support Intern', 'Assist with troubleshooting hardware and software issues. 3-month paid internship.', 'networking,troubleshooting,communication', 'KES 20,000', 'Mombasa', 'open'),
(2, 'Farm Operations Assistant', 'Support daily farm operations including record keeping and logistics coordination.', 'logistics,record-keeping,communication', 'KES 30,000 - 40,000', 'Nakuru', 'open'),
(2, 'Social Media Manager', 'Manage our brand presence across social platforms. Content creation and analytics.', 'social-media,content-creation,graphic-design', 'KES 35,000 - 55,000', 'Remote', 'open');

-- Seed Applications
INSERT INTO applications (job_id, jobseeker_id, status) VALUES
(1, 3, 'pending'),
(2, 4, 'accepted'),
(5, 5, 'pending');
