# Hướng dẫn Cài đặt Docker và Chạy Dự án "Bambu Farm" trên Windows

Tài liệu này hướng dẫn chi tiết từng bước (step-by-step) để bạn có thể tự cài đặt Docker và chạy toàn bộ hệ thống Bambu Farm trên máy tính Windows.

---

## Phần 1: Cài đặt Docker Desktop trên Windows

### Bước 1: Kiểm tra yêu cầu hệ thống
Để cài được Docker, máy tính Windows của bạn cần thoả mãn hai điều kiện sau:
1. Đang sử dụng **Windows 10 (64-bit)** bản Pro/Enterprise/Education hoặc **Windows 11**.
2. Đã bật tính năng **Virtualization (Ảo hóa)** trong BIOS (thường đa số máy tính hiện nay đều đã bật sẵn). Để chắc chắn, bạn có thể bấm tổ hợp phím `Ctrl + Shift + Esc` để mở **Task Manager**, chuyển sang tab **Performance**, chọn **CPU** và nhìn ở góc dưới bên phải xem dòng **Virtualization** có đang là `Enabled` không.

### Bước 2: Bật tính năng WSL 2 (Windows Subsystem for Linux)
Docker trên Windows chạy mượt mà nhất khi dùng WSL 2.
1. Mở **Start Menu**, gõ tìm kiếm chữ `cmd` hoặc `PowerShell`.
2. Nhấp chuột phải vào biểu tượng **Command Prompt** (hoặc PowerShell) và chọn **Run as administrator**.
3. Copy và dán dòng lệnh sau vào màn hình đen rồi nhấn Enter:
   ```cmd
   wsl --install
   ```
4. Đợi quá trình tải và cài đặt hoàn tất. **Khởi động lại máy tính** để áp dụng thay đổi.

### Bước 3: Tải và Cài đặt Docker Desktop
1. Truy cập trang chủ Docker: [Tải Docker Desktop cho Windows](https://www.docker.com/products/docker-desktop/).
2. Nhấp vào nút **Download for Windows** để tải file cài đặt `.exe` về máy.
3. Mở file vừa tải về (`Docker Desktop Installer.exe`).
4. Trong màn hình cài đặt, đảm bảo đã **tích chọn** dòng **"Use WSL 2 instead of Hyper-V"**.
5. Bấm **Ok** và chờ phần mềm cài đặt. Khi cài xong, bấm **Close and restart** (nếu máy yêu cầu).

### Bước 4: Mở và cấu hình Docker
1. Cài xong, tìm biểu tượng con cá voi màu xanh (**Docker Desktop**) trên màn hình Desktop hoặc Start Menu để mở lên.
2. Lần đầu mở, phần mềm sẽ yêu cầu bạn chấp nhận điều khoản (Accept Terms), hãy nhấn **Accept**.
3. Cửa sổ Docker hiện ra, nếu thấy biểu tượng ở góc dưới bên trái có màu xanh lá cây và báo chữ "Engine running", tức là Docker đã sẵn sàng hoạt động.
*(Lưu ý: Bạn có thể chọn "Skip" nếu Docker hỏi khảo sát người dùng đầu vào)*.

---

## Phần 2: Khởi chạy dự án Bambu Farm

### Bước 1: Chuẩn bị mã nguồn
Mở thư mục chứa mã nguồn của bạn (thư mục `Bambu-farm`). Đảm bảo rằng trong thư mục này có file tên là `docker-compose.yml`.

### Bước 2: Thiết lập biến môi trường
Dự án đã có sẵn các file `.env.example` ở cả phần Giao diện (frontend) và Máy chủ (backend). Bạn cần tạo file cấu hình thực tế cho chúng:

1. Vào thư mục `backend/`, copy file `.env.example` và đổi tên bản copy thành `.env`. Trong file `.env` này đã có sẵn các cấu hình cổng (port), PostgreSQL, Redis cần thiết để docker-compose tự hoạt động.
2. Vào thư mục `frontend/`, copy file `.env.example` và đổi tên thành `.env.production` (hoặc `.env.local` nếu bạn chạy dev).
*(Nếu bạn không sửa lại cấu hình máy chủ, bạn có thể để nguyên mọi thông số mặc định).*

### Bước 3: Mở Terminal (Command Prompt) tại thư mục dự án
1. Chọn thanh địa chỉ (Address bar) trong cửa sổ quản lý thư mục đang mở mã nguồn `Bambu-farm`.
2. Xóa tất cả chữ trên thanh địa chỉ, Gõ chữ `cmd` và nhấn Enter. 
3. Một cửa sổ lệnh màu đen sẽ được bật lên, và vị trí hiện tại đã trỏ đúng vào thư mục code của bạn.

### Bước 4: Xây dựng và khởi chạy các dịch vụ
Tại cửa sổ dòng lệnh ở Bước 3, gõ lệnh sau và nhấn Enter:

```cmd
docker-compose up --build -d
```

**Giải thích ý nghĩa của lệnh:**
- `docker-compose up`: Yêu cầu Docker đọc file `docker-compose.yml` để tạo các container.
- `--build`: Bắt buộc Docker phải đóng gói lại (build) phần **frontend** và **backend** mới nhất theo file `Dockerfile` của bạn.
- `-d` (detach): Chạy ngầm các chương trình ở chế độ ẩn, giúp bạn vẫn có thể gõ tiếp lệnh khác trên màn hình này.

*Lưu ý: Quá trình chạy lệnh này lần đầu tiên có thể mất từ 5-10 phút tuỳ vào tốc độ mạng, vì Docker phải tải hệ điều hành giả lập, tải Node.js (cho Frontend) và tải Golang (cho Backend) về máy bạn.*

### Bước 5: Kiểm tra trạng thái hoạt động
Sau khi quá trình trên chạy xong báo chữ `Started` hoặc `Running` toàn bộ các dòng. Bạn hãy kiểm tra lại bằng cách gõ:

```cmd
docker-compose ps
```

Nếu bạn thấy danh sách các hệ thống như `bambu-frontend`, `bambu-backend`, `bambu-db` (Postgres) và `bambu-redis` đều có trạng thái là `Up`, tức là dự án của bạn đã chạy thành công rực rỡ!

### Bước 6: Truy cập ứng dụng
Mở trình duyệt web của bạn (Chrome, Edge, Cốc Cốc...) và truy cập vào các đường dẫn sau:

- **Giao diện người dùng (Frontend):** `http://localhost:3000`
- **Máy chủ dữ liệu (Backend API):** `http://localhost:8080` (Bạn có thể kiểm tra sức khoẻ server qua `http://localhost:8080/health`)

### Bước 7: Tự động khởi động cùng Windows (Tùy chọn)
Nếu bạn muốn hệ thống tự động chạy mỗi khi bật máy:
1. Chạy file `install_startup.bat`.
2. Từ giờ, Bambu Farm sẽ tự khởi động ngầm mỗi khi bạn mở máy tính.

---

## Các công cụ (Scripts) đi kèm

Chúng tôi cung cấp một số công cụ để bạn thao tác nhanh hơn:

- **`start_hidden.vbs`**: Chạy toàn bộ hệ thống dưới nền ẩn (không hiện cửa sổ đen CMD).
- **`stop_app.bat`**: **Nút Tắt** nhanh toàn bộ hệ thống.
- **`install_startup.bat`**: Cài đặt để hệ thống tự chạy khi bật máy.
- **`remove_startup.bat`**: Gỡ bỏ tính năng tự chạy khi bật máy.

---

## Các thao tác (Commands) hữu ích nên biết

Nếu bạn muốn **dừng** và tắt toàn bộ hệ thống:
```cmd
docker-compose down
```

Nếu bạn muốn **xem log (lịch sử lỗi/hoạt động)** của máy chủ Backend:
```cmd
docker-compose logs -f backend
```

Nếu bạn sửa code và muốn hệ thống **cập nhật lại mã nguồn mới**:
```cmd
docker-compose up --build -d
```

---
**Tác giả:** Quang Huy with Antigravity
