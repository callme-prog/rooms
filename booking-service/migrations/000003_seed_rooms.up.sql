INSERT INTO rooms (id, name, description, price, capacity, image_url) VALUES
  (gen_random_uuid(), 'Стандартный номер', 'Уютный номер с видом на двор, двуспальная кровать, Wi-Fi',    4500.00, 2, 'https://images.unsplash.com/photo-1631049307264-da0ec9d70304?w=800'),
  (gen_random_uuid(), 'Улучшенный номер', 'Просторный номер с балконом и панорамным видом на город',     6800.00, 2, 'https://images.unsplash.com/photo-1590490360182-c33d57733427?w=800'),
  (gen_random_uuid(), 'Семейный номер',   'Два раздельных спальных места, идеально для семей с детьми',  8500.00, 4, 'https://images.unsplash.com/photo-1566665797739-1674de7a421a?w=800'),
  (gen_random_uuid(), 'Люкс',             'Апартаменты класса люкс: гостиная, спальня, джакузи',        14000.00, 2, 'https://images.unsplash.com/photo-1582719478250-c89cae4dc85b?w=800'),
  (gen_random_uuid(), 'Президентский люкс','Два этажа, терраса с видом на море, персональный butler', 32000.00, 4, 'https://images.unsplash.com/photo-1618773928121-c32242e63f39?w=800')
ON CONFLICT DO NOTHING;

CREATE INDEX IF NOT EXISTS idx_bookings_user_id   ON bookings(user_id);
CREATE INDEX IF NOT EXISTS idx_bookings_room_id   ON bookings(room_id);
CREATE INDEX IF NOT EXISTS idx_bookings_status    ON bookings(status);
CREATE INDEX IF NOT EXISTS idx_bookings_check_in  ON bookings(check_in);
