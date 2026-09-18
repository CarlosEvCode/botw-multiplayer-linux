# The Legend of Zelda: Breath of the Wild Multiplayer en Linux / SteamOS

[English (README.md)](README.md) | [Español (README.es.md)](README.es.md)

Instalador automatizado, TUI Manager interactivo y entorno técnico completo para hostear y jugar a *The Legend of Zelda: Breath of the Wild Multiplayer* (*Milk Bar Launcher* / *Cemu 1.26.2*) en distribuciones GNU/Linux y SteamOS mediante Wine.

---

## 1. Instalación Rápida

Ejecuta este comando en tu terminal para descargar, configurar e instalar todo automáticamente:

```bash
curl -sSL https://raw.githubusercontent.com/CarlosEvCode/botw-multiplayer-linux/main/install.sh | bash
```

### Instalación Manual (Vía Git)

```bash
git clone https://github.com/CarlosEvCode/botw-multiplayer-linux.git
cd botw-multiplayer-linux
chmod +x install.sh
./install.sh
```

> **Nota:** El instalador aprovisiona automáticamente Wine de 64 bits, Microsoft .NET 8.0 Desktop Runtime, Visual C++ 2015–2022, Cemu 1.26.2 (con los parches BCML y los 32 modelos de Link ya fusionados), Milk Bar Launcher y la interfaz CLI `botw-manager`.

---

## 2. TUI Manager (`botw-manager`)

Administra servidores, clientes, modos de juego, rutas del juego y direcciones IP de red directamente desde tu terminal:

```bash
botw-manager
```

### Funciones y Atajos de Teclado
* **Navegación Global:** Presiona `Tab` / `Shift+Tab` para alternar entre pestañas en cualquier momento.
* **1. Dashboard:**
  * `[s]`: Iniciar / Detener Servidor Dedicado (con consola de logs en vivo).
  * `[m]`: Lanzar Milk Bar Launcher.
  * `[c]`: Lanzar Cemu 1.26.2.
  * `[t]`: Copiar IP de Tailscale / LAN al portapapeles.
  * `[q]`: Salir de la aplicación.
* **2. Gamemodes:**
  * `←` / `→`: Alternar modos de juego:
    * **Cooperativo Estándar (Libre):** Modo cooperativo con sincronizaciones totalmente personalizables.
    * **Hunter vs Speedrunner:** Modo competitivo Manhunt (bloquea el progreso individual automáticamente).
    * **DeathSwap:** Modo de supervivencia con intercambio de posiciones periódicas (bloquea sincronizaciones de progreso).
  * `Espacio`: Alternar reglas de sincronización individuales (Misiones, Santuarios, Torres, Kologs, Enemigos, Mazmorras, Ubicaciones).
  * `s`: Guardar reglas (reinicia el servidor automáticamente si está corriendo).
* **3. Rutas & DLC:**
  * `↓` / `↑`: Seleccionar campo (Juego Base, Update v208, DLC v80).
  * `f`: Abrir el explorador de carpetas gráfico del sistema (*Zenity* / *Kdialog*) para elegir rutas con un clic.
  * `Enter`: Guardar y reconstruir los enlaces simbólicos de Wine y BCML.
* **4. Red:**
  * Muestra las interfaces activas de Tailscale VPN, Red Local LAN y ZeroTier con instrucciones de conexión.
* **5. Ajustes:**
  * **Selector de Idioma:** Cambia instantáneamente entre **English** y **Español**.
  * Información del proyecto, versión y enlaces al repositorio.

---

## 3. Guía de Conexión Multijugador

### Anfitrión (Host / Servidor)
1. Abre `botw-manager` y presiona `[s]` para iniciar el Servidor Dedicado.
2. Presiona `[m]` para abrir Milk Bar Launcher.
3. En Milk Bar Launcher, ingresa:
   * **IP:** `127.0.0.1`
   * **Puerto:** `5050`
4. Haz clic en **Connect** (Cemu se abrirá automáticamente sincronizado).

### Jugadores Remotos (Amigos por Internet)
1. Instala [Tailscale](https://tailscale.com/) (o ZeroTier) tanto en el equipo del Host como en el de los clientes.
2. El Host copia su IP de Tailscale (presionando `[t]` en `botw-manager`).
3. Los jugadores remotos abren Milk Bar Launcher e ingresan:
   * **IP:** `IP de Tailscale del Host` (ej. `100.x.y.z`)
   * **Puerto:** `5050`
4. Hacen clic en **Connect**.

---

## 4. Estabilidad de Ventanas en Hyprland / Wayland

Si usas **Hyprland**, añade las siguientes reglas a tu archivo `~/.config/hypr/hyprland.lua` para evitar parpadeos en XWayland y enviar Milk Bar Launcher en segundo plano al Workspace 3:

```lua
o.window({ class = ".*(milk bar launcher|MilkBar).*" }, {
  float = true,
  center = true,
  size = "1188 670",
  suppress_event = "activate maximize fullscreen",
  workspace = "3 silent",
})
o.window({ class = "^(cemu\\.exe)$" }, {
  opaque = true,
  suppress_event = "maximize",
})
```

---

## 5. Arquitectura Técnica y Notas de Compatibilidad

Como referencia técnica, se implementaron las siguientes soluciones para garantizar la estabilidad bajo Wine en Linux:

1. **Wineprefix Compartido de 64 bits (`~/.local/share/wineprefixes/botw-multiplayer`):**
   * Cemu y Milk Bar Launcher deben ejecutarse en el mismo prefijo para que `InjectDLL.dll` y los Named Pipes puedan conectarse a la memoria de Cemu sin bloqueos de pantalla blanca.
2. **Microsoft .NET 8.0 Desktop Runtime x64:**
   * Instalado directamente en el Wineprefix para solucionar cierres con código de salida 150.
3. **Validación de Paquetes de Modelos BCML:**
   * Incluye los 32 paquetes de modelos de Link (`Jugador1.sbactorpack` a `Jugador32.sbactorpack`) en `store_dir/merged/content/Actor/Pack/` enlazados a `BreathOfTheWild_BCML`.
4. **Enlace de Socket del Servidor:**
   * `ServerConfig.ini` se configura con `IP=127.0.0.1` y `DefaultGamemode=True` para enlazar inmediatamente el puerto `5050` y evitar excepciones de resolución DNS de `localhost`.
5. **Sanitización de Recursos en AppData Roaming:**
   * Extrae `QuestFlagsNames.txt` y `ArmorMapping.txt` en `%APPDATA%/BOTWM/` para evitar errores `FileNotFoundException`.

---

## 6. Créditos y Atribuciones

* **Milk Bar Launcher & Dedicated Server:** Desarrollado por la comunidad de [MilkBarModding](https://github.com/MilkBarModding/MilkBarLauncher).
* **Cemu (Emulador de Wii U):** Desarrollado por [Team Cemu](https://cemu.info/) bajo licencia Mozilla Public License 2.0.
* **BCML:** Creado por [NiceneNerd](https://github.com/NiceneNerd/BCML).
* **Linux Manager y Despliegue:** Mantenido por [CarlosEvCode](https://github.com/CarlosEvCode).

### Descargo de Responsabilidad
Este repositorio no contiene ni distribuye archivos protegidos por derechos de autor de Nintendo, ejecutables oficiales (`.rpx`), claves de cifrado ni contenido propietario. Los usuarios deben suministrar su propia copia legal de *The Legend of Zelda: Breath of the Wild*, sus actualizaciones y DLC.
