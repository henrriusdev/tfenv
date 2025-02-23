# tfenv - Terraform Environment Manager

**tfenv** is a simple CLI tool that converts `.env` files into `terraform.tfvars` for easy integration with Terraform. It supports Windows, Linux, and macOS.

---

## 🚀 Installation

### **Go**

You can install this by using
```bash
go install github.com/henrriusdev/tfenv@latest
```

### **Windows**

Run the following command in PowerShell:

```powershell
Invoke-WebRequest -Uri "https://github.com/YOUR-USER/tfenv/releases/latest/download/install.ps1" -OutFile "install.ps1"
.\install.ps1
```

After installation, restart the terminal and run:

```powershell
tfenv
```

---

### **Linux/macOS**

Run the following command in Terminal:

```sh
curl -fsSL https://github.com/YOUR-USER/tfenv/releases/latest/download/install.sh | bash
```

Once installed, you can run:

```sh
tfenv
```

---

## 🎯 Usage

1. Navigate to your project directory.
2. If the `.env` file exists, just run:
   ```sh
   tfenv
   ```
   Otherwise, you will be prompted to enter the `.env` file path.
3. The tool will convert it to `terraform.tfvars` and optionally generate `variables.tf`.

---

## 🛠 Features

- **Converts .env to terraform.tfvars**
- **Generates variables.tf (optional)**
- **Interactive CLI** with [Bubble Tea](https://github.com/charmbracelet/bubbletea) and [Huh](https://github.com/charmbracelet/huh)
- **Works on Windows, Linux, and macOS**
- **Can be installed with Go**

---

## 🔧 Contributing

1. Clone the repository:
   ```sh
   git clone https://github.com/YOUR-USER/tfenv.git
   ```
2. Navigate to the project:
   ```sh
   cd tfenv
   ```
3. Build the project:
   ```sh
   go build -o tfenv
   ```
4. Run it:
   ```sh
   ./tfenv
   ```

---

## 📜 License

This project is licensed under the MIT License.

