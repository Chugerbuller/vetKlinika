# Vet Clinic Management System (.NET Razor Pages Rewrite)

  ## 📋 Project Overview
  - **Domain:** Animal hospital / veterinary clinic management system
  - **Goal:** Rewrite existing application to modern .NET using Razor Pages
  - **Working Directory:** `D:\code\vetKlinika`

  ## 🛠 Tech Stack & Conventions
  - **Framework:** .NET 8+, C# 12, ASP.NET Core Razor Pages
  - **Database:** EF Core + sqlite
  - **Auth:** ASP.NET Identity (default), switch to JWT/external providers later if needed
  - **Frontend:** Razor Views + Tag Helpers, Bootstrap/Tailwind, HTMX/AJAX for lightweight interactivity
  - **Architecture:** Clean-ish layering: `Pages/` → `Services/` → `Data/` (EF Core) → `Models/`
  - **DI & Config:** Built-in ASP.NET DI, `appsettings.json`, Options pattern for config
  - **Code Style:** Nullable enabled, implicit usings, standard C# conventions, XML comments on public APIs

  ## 🐾 Expected Core Domains
  | Domain | Key Properties / Notes |
  |--------|------------------------|
  | `Patient` (Animal) | Name, species, breed, DOB, microchip, photo, status (active/inactive) |
  | `Owner` | Contact info, address, payment method, linked patients |
  | `Appointment` | Date/time, assigned vet, type/status, notes, patient reference |
  | `Veterinarian` | Staff profile, specialties, availability, contact |
  | `Service/Treatment` | Catalog name, duration,