# Bambu Farm Manager

Bambu Farm Manager là một hệ thống quản lý trang trại máy in 3D BambuLab tập trung, cho phép theo dõi và điều khiển nhiều máy in cùng lúc qua mạng LAN.

**Dự án được thực hiện bởi: Quang Huy with Antigravity**

## Tính năng chính
- **Quản lý máy in:** Thêm, sửa, xóa các máy in trong hệ thống.
- **Tự động tìm kiếm:** Tự động phát hiện các máy in BambuLab trong mạng LAN qua mDNS.
- **Theo dõi thời gian thực:** Xem nhiệt độ đầu in, bàn nhiệt và tiến trình in qua WebSockets.
- **Quản lý lệnh in:** Gửi file, tạm dừng, tiếp tục hoặc hủy lệnh in.
- **Hệ thống cảnh báo:** Gửi thông báo lỗi hoặc bất thường qua Email và Telegram.
- **Luồng Camera:** Xem trực tiếp hình ảnh từ máy in qua giao diện web.
- **Phân quyền người dùng:** Đăng nhập và quản lý theo tổ chức (Organization) với JWT.

## Công nghệ sử dụng
- **Backend:** Golang 1.21+ (Gin, GORM, MQTT, Redis, WebSockets).
- **Frontend:** Next.js 14+ (App Router, TailwindCSS).
- **Cơ sở dữ liệu:** PostgreSQL.
- **Hàng đợi:** Redis.
- **Đóng gói:** Docker & Docker Compose.

## Hướng dẫn cài đặt
Vui lòng xem file [SETUP_DOCKER_WIN.md](./SETUP_DOCKER_WIN.md) để biết chi tiết cách cài đặt trên Windows.

## Giấy phép
Dự án được phát hành dưới giấy phép MIT.
