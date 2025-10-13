-- General
UPDATE permissions SET group_name = 'General' WHERE name LIKE 'general.%';

-- User
UPDATE permissions SET group_name = 'User' WHERE name LIKE 'user.%';

-- Role
UPDATE permissions SET group_name = 'Role' WHERE name LIKE 'role.%';

-- Permission
UPDATE permissions SET group_name = 'Permission' WHERE name LIKE 'permission.%';

-- Product Stats
UPDATE permissions SET group_name = 'Product Stats' WHERE name LIKE 'product_stats.%';

-- Category
UPDATE permissions SET group_name = 'Category' WHERE name LIKE 'category.%';

-- Address
UPDATE permissions SET group_name = 'Address' WHERE name LIKE 'address.%';

-- Review
UPDATE permissions SET group_name = 'Review' WHERE name LIKE 'review.%';

-- Wishlist
UPDATE permissions SET group_name = 'Wishlist' WHERE name LIKE 'wishlist.%';

-- Brand
UPDATE permissions SET group_name = 'Brand' WHERE name LIKE 'brand.%';

-- Color
UPDATE permissions SET group_name = 'Color' WHERE name LIKE 'color.%';

-- Attribute Type
UPDATE permissions SET group_name = 'Attribute Type' WHERE name LIKE 'attribute_type.%';

-- Attribute Option
UPDATE permissions SET group_name = 'Attribute Option' WHERE name LIKE 'attribute_option.%';

-- Size Category
UPDATE permissions SET group_name = 'Size Category' WHERE name LIKE 'size_category.%';

-- Size Option
UPDATE permissions SET group_name = 'Size Option' WHERE name LIKE 'size_option.%';
