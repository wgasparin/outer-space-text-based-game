# **Game Design Document (GDD)** 



# Module 1: The Physics (Attributes & State)

These are the core variables that define every entity (Player or NPC) in your database. All interactions in the game are mathematical operations on these four values.

### 1.1 Core Stats (The "Exosuit Telemetry")

| Stat Name | Analog | Description | Backend Implication |
| :--- | :--- | :--- | :--- |
| **Integrity** | HP  | Structural health. If 0, the unit is destroyed. | Stored as `current_integrity` and `max_integrity`. |
| **Firepower** | Strength | Raw weapon output before mitigation. | Used in the Damage Formula. |
| **Shielding** | Defense | Energy field that reduces incoming damage. | Used in the Damage Formula. |
| **Thrusters** | Speed | System speed and maneuverability. | Determines the **Global Cooldown (GCD)** and **Warp Speed**. |

### 1.2 Secondary Resources

  * **Reactor Core (Energy):**
      * **Max:** 100 (Fixed for all classes).
      * **Regen:** 5 Energy per second.
      * **Usage:** Required to cast "Tactical Modules" (Special skills).
      * *Dev Note:* Do not update this every second in the DB. Calculate it on demand: `Current = Stored + ((Now - LastUpdate) * Rate)`.

-----

# Module 2: The Loadouts (Character Classes)

When a player creates a character (`POST /players`), they select one Loadout. This applies base modifiers and defines their available methods (skills).

**Base Stats for all:** Integrity: 100, Firepower: 10, Shielding: 10, Thrusters: 100.

### 2.1 Class Definitions

| Class | Role | Stat Modifiers | Passive Trait (The "If" Logic) |
| :--- | :--- | :--- | :--- |
| **Striker** | DPS | Firepower +50%<br>Integrity -20% | **Overheat:** Damage dealt increases by 1% for every 1% of missing Integrity. |
| **Titan** | Tank | Shielding +50%<br>Thrusters -20% | **Ablative Armor:** Incoming damage is capped at 10% of Max Integrity per hit. |
| **Spectre** | Assassin | Firepower +20%<br>Thrusters +30% | **Void Cloak:** First attack on a target not targeting you ignores 100% Shielding. |
| **Engineer** | Healer | Thrusters +20%<br>Firepower -20% | **Nanites:** Healing above Max Integrity creates a temporary Barrier (Max 20% Integrity). |
| **Navigator** | Support | Integrity +20%<br>Shielding +20% | **Telemetry Link:** Buffs applied to allies have 50% longer duration. |

### 2.2 Skill Sets (The Methods)

Every class has two core abilities: a spammable Primary Weapon and a powerful, resource-intensive Tactical Module. These are the `methods` that can be called on a character object.

*   **Primary Weapon:** 0 Energy Cost. The bread-and-butter attack. Bound to the Global Cooldown (GCD).
*   **Tactical Module:** High Energy Cost. A high-impact ability that defines the class's role. Also bound to the GCD.

| Class      | Primary Weapon                                  | Tactical Module (Cost: Energy)                                 | Skill Description                                                                                                                              |
| :---       | :---                                            | :---                                                           | :---                                                                                                                                           |
| **Striker**  | **Plasma Bolt:** 100% Firepower Dmg.            | **Cluster Bomb (40):** 3x attacks at 60% Firepower instantly.    | **Plasma Bolt:** A straightforward, reliable damage beam. <br> **Cluster Bomb:** Unleashes a rapid volley, perfect for bursting down a single target. The multiple hits are resolved sequentially. |
| **Titan**    | **Suppressing Fire:** 80% Firepower + High Threat. | **Energy Barrier (50):** +100% Shielding for 5s. Rooted in place. | **Suppressing Fire:** Generates high "threat," making NPCs more likely to attack the Titan. <br> **Energy Barrier:** Temporarily transforms the Titan into an immovable fortress, absorbing immense damage but preventing any movement. |
| **Spectre**   | **Void Blade:** 100% Firepower Dmg.             | **Phase Shift (30):** Warp to random adjacent sector (Escape). | **Void Blade:** A standard melee-range attack. <br> **Phase Shift:** An emergency escape maneuver. The server selects a random valid `Exit` from the current sector and initiates a `Warp` state, effectively allowing the Spectre to flee from a losing fight. |
| **Engineer** | **Welding Beam:** 50% Firepower Dmg.            | **Reconstruction (40):** Heal target for $3 \times Firepower$.  | **Welding Beam:** A weak offensive beam that doubles as a field tool. <br> **Reconstruction:** A powerful single-target heal. The healing amount scales with the Engineer's own Firepower, creating an interesting gear choice. Can be self-cast. |
| **Navigator**| **Target Lock:** 80% Firepower Dmg.             | **Full Impulse (60):** All allies in sector get +50 Thrusters for 10s. | **Target Lock:** A standard attack that also "paints" the target for allies. <br> **Full Impulse:** A fleet-wide buff that significantly increases the movement speed and action rate (lower GCD) of all friendly players in the same sector. Essential for coordinated fleet actions. |

-----

# Module 3: The Starchart (Map & Movement)

The world is a **Directed Graph**.

  * **Nodes:** Sectors (Star Systems).
  * **Edges:** Warp Lanes.

### 3.1 Sector Properties

Each Sector Object in the database contains:

  * `ID`: Unique UUID.
  * `Type`: **Safe** (No combat), **Open** (PVE/PVP), or **Hazard** (Environmental damage).
  * `Exits`: Array of adjacent Sector IDs (e.g., `["sector_02", "sector_05"]`).

### 3.2 The Warp Mechanics (Movement Rules)

Movement is **not instant**. It is a state change that takes time.

1.  **Initiation:** Player requests `POST /move { target: "sector_02" }`.
2.  **Validation:** Is `sector_02` in the `Exits` list of the current sector?
3.  **Calculation:**
    $$WarpTime_{ms} = \frac{5000}{(Thrusters / 100)}$$
    *(Base time is 5 seconds. Higher thrusters = faster travel).*
4.  **State Change:** Player state becomes `WARPING`. They cannot attack.
5.  **Completion:** When the timer expires, update `player.currentSector` and set state to `IDLE`.

-----

# Module 4: Combat Protocols (Engagement Rules)

This is the core loop for your `POST /attack` endpoint.

### 4.1 The Global Cooldown (GCD)

To prevent script spamming and manage server load, actions are time-gated.
$$GCD_{ms} = \frac{2000}{(Thrusters / 100)}$$

  * *Rule:* After any action, the player receives a `busyUntil` timestamp. Any request received before this timestamp returns `429 Too Many Requests`.

### 4.2 The Damage Formula

Combat is deterministic (no dice).
$$Damage = max(1, (AttackerFirepower \times SkillMultiplier) - TargetShielding)$$

  * *Note:* Minimum damage is always 1.

### 4.3 Engagement States

A player can be in one of these mutually exclusive states:

1.  **IDLE:** Ready to act. Energy regenerating.
2.  **COMBAT:** Engaged. Energy regenerating. Cannot `Warp` immediately (must wait 10s after last combat action).
3.  **WARPING:** Moving between sectors. Cannot Attack.
4.  **DEAD:** Integrity = 0. Must Respawn (`POST /respawn`).

### 4.4 Targeting Logic

  * **Range:** You can only attack targets in the **same Sector**.
  * **Friendly Fire:**
      * In **Safe Sectors**: Damage is 0.
      * In **Open Sectors**: You can attack anyone (if the game is PVP).

-----

# Module 5: Progression (Scaling)

How players grow stronger over time.

### 5.1 Experience (XP)

  * **Gain:** Killing a target grants XP equal to the target's `MaxIntegrity`.
  * **Level Up Threshold:**
    $$XP_{Required} = CurrentLevel^2 \times 100$$

### 5.2 Level Up Rewards

Upon reaching a new level:

1.  Integrity fully restored.
2.  **System Points:** Player receives **5 Points** to manually distribute into Integrity, Firepower, Shielding, or Thrusters via `POST /upgrade`.

-----

### Developer Implementation Tips

1.  **Start with Module 1 & 2:** Create the `Player` class/model and the database schema.
2.  **Build Module 3:** Create the map graph (JSON or DB) and the `move` endpoint.
3.  **Build Module 4:** This is the hardest part. Implement the `attack` endpoint with the GCD logic.
4.  **Race Conditions:** When implementing combat, ensure two players attacking the same target don't reduce its HP below zero twice. Use **Atomic Operations** (e.g., `UPDATE enemies SET integrity = integrity - 10 WHERE id = 1 AND integrity > 0`).