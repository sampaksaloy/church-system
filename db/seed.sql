-- Church System Sample Data
-- Run this after initial setup to populate with sample records

-- Sample Announcements
INSERT INTO announcements (title, content, category, is_pinned, is_active, admin_id) VALUES
('Welcome to Our Parish Website!', 'We are excited to launch our new church website where you can find all the latest announcements, event schedules, and connect with our community. Check back regularly for updates!', 'General', true, true, 1),
('Sunday Mass Schedule Change', 'Please be informed that starting next month, the Sunday morning Mass will be moved from 8:00 AM to 7:30 AM to accommodate more parishioners. The 10:30 AM Mass schedule remains unchanged.', 'Mass Schedule', false, true, 1),
('Youth Fellowship This Saturday', 'The Youth Ministry is holding its monthly fellowship this Saturday at 5:00 PM in the Parish Hall. All youth aged 13-25 are welcome. Bring your friends! We will have games, sharing, and a special talk on faith.', 'Youth', false, true, 1),
('Choir Practice Schedule', 'The Parish Choir invites all interested members to join our weekly practice sessions every Wednesday at 6:30 PM. New members are welcome! No prior experience required, just a heart for worship.', 'Choir', false, true, 1),
('Community Feeding Program', 'Our parish will be conducting a community feeding program for underprivileged families this coming Sunday after the 10:30 AM Mass. Donations of cooked food, bottled water, and hygiene kits are greatly appreciated.', 'Community', false, true, 1);

-- Sample Events
INSERT INTO events (title, description, location, event_date, start_time, end_time, category, is_recurring, is_active, admin_id) VALUES
('Sunday Mass', 'Regular Sunday morning Mass. All parishioners are welcome.', 'Main Church', CURRENT_DATE + INTERVAL '3 days', '08:00', '09:00', 'Mass', true, true, 1),
('Sunday Mass (2nd Schedule)', 'Second Sunday Mass schedule for working parishioners.', 'Main Church', CURRENT_DATE + INTERVAL '3 days', '10:30', '11:45', 'Mass', true, true, 1),
('Youth Fellowship Meeting', 'Monthly youth gathering with games, sharing, and faith formation talk.', 'Parish Hall', CURRENT_DATE + INTERVAL '5 days', '17:00', '20:00', 'Youth', false, true, 1),
('Bible Study Group', 'Weekly Bible study and reflection for all parishioners.', 'Parish Center Room 2', CURRENT_DATE + INTERVAL '7 days', '19:00', '20:30', 'Bible Study', true, true, 1),
('Choir Practice', 'Weekly choir rehearsal session for all choir members.', 'Church Choir Loft', CURRENT_DATE + INTERVAL '9 days', '18:30', '20:00', 'Choir', true, true, 1),
('Community Outreach Program', 'Parish feeding program and distribution of goods to underprivileged families in the community.', 'Covered Court, Parish Grounds', CURRENT_DATE + INTERVAL '10 days', '13:00', '17:00', 'Community Outreach', false, true, 1),
('Parish Pastoral Council Meeting', 'Monthly meeting of the Parish Pastoral Council to discuss parish affairs and upcoming activities.', 'Parish Office', CURRENT_DATE + INTERVAL '14 days', '09:00', '11:00', 'General', false, true, 1);
