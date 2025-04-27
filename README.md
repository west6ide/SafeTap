# SafeTap

SafeTap is a mobile security solution designed to provide real-time protection, emergency alerts, and location sharing capabilities. It combines a powerful Golang backend with mobile client applications to help users stay safe in critical situations.

## 🔍 Overview
SafeTap offers users an immediate way to send an SOS signal, share their real-time location, manage emergency contacts, and trigger fake calls for discreet safety maneuvers. The backend system handles secure authentication, location updates, and notification delivery.

## 🌐 Features
- 📲 **Quick SOS Activation**: Instantly alert emergency contacts with your current location.
- 📊 **Live Location Sharing**: Update and broadcast your real-time GPS location.
- 📢 **Emergency Notifications**: Send automatic push notifications.
- ☎️ **Direct Emergency Calling**: Integrate direct calls to emergency services (102).
- 🤖 **Fake Call Simulation**: Set up fake calls to create safe exit strategies.
- 👥 **Trusted Emergency Contacts**: Manage personal trusted contact lists.
- 📈 **Dangerous Person Alerts**: Mark and view dangerous individuals in your area.
- 🔑 **Secure Google OAuth Authentication**: Fast and secure login.

## 👨‍💼 Technologies Used
- **Backend**: Golang (net/http)
- **Database**: PostgreSQL
- **Authentication**: Google OAuth 2.0
- **Push Notifications**: Firebase Cloud Messaging (FCM)
- **Location Services**: Google Maps API
- **Mobile App**: Android (Kotlin)

## 🌐 Project Structure
```
SafeTap/
├── config/                # Database connection setup
├── controller/            # All API route controllers
│   ├── authentication/   # Google authentication logic
├── users/                 # Data models (User, FakeCall, DangerousPerson)
├── .idea/                 # IDE config files
└── .git/                  # Git repository files
```

## 📆 Installation

### Backend Setup
```bash
# Clone the repository
$ git clone https://github.com/west6ide/SafeTap.git
$ cd SafeTap

# Configure your PostgreSQL database credentials inside config/db.go

# Run the server
$ go run main.go
```

### Android Mobile App
- Open the Android project (not included in this repository) in Android Studio.
- Configure the backend API URL.
- Build and run the app on a device.

## 🛋️ API Endpoints
- **POST** `/sos/send` - Send SOS signal.
- **POST** `/location/update` - Update user location.
- **GET** `/notifications` - Get all notifications.
- **POST** `/fake-call/create` - Schedule a fake call.
- **POST** `/contacts/add` - Add an emergency contact.
- **POST** `/auth/google` - Google sign-in.

## 🔄 Main Modules
| Module | Description |
| :--- | :--- |
| `sosFunc.go` | Handling SOS signal sending |
| `updateLocation.go` | Real-time location updates |
| `notificationSos.go` | Sending notifications |
| `fakeCallController.go` | Managing fake calls |
| `auth.go` | Basic authentication |
| `google_auth.go` | Google OAuth2 authentication |
| `user.go`, `dangerousPerson.go`, etc. | Database models |

## 💡 Future Improvements
- iOS application
- Two-way chat with emergency contacts
- Danger zone mapping and alerts
- AI-based emergency prediction

## 🌟 Contributing
Pull requests are welcome! For major changes, please open an issue first to discuss what you would like to change.

## 👁️ License
This project is licensed under the [MIT License](LICENSE).

---

> SafeTap — Your real-time guardian. Stay connected. Stay protected.

