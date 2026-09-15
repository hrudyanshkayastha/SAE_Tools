# SAE TheHive Source Audit

## 1. Upstream Repository
- **Remote**: `https://github.com/TheHive-Project/TheHive.git`
- **Commit**: `d390a031c6a2e4e049969623e160a0a55e2dbd73`
- **Version**: TheHive 4 / 5 (Branch/Head)

## 2. Source Integrity
- The local working tree was heavily modified and incorrectly branded as `sae case management`.
- All `thehive` paths inside `app/org/thp/` and `client/` were renamed to `sae_case_management`.
- These structural corruptions were fully reversed using `git reset --hard` and `git clean -fd`.
- The pristine upstream source code was safely relocated to `SAE/thehive`.

## 3. SAE Integration
- SAE maintains the pristine upstream TheHive source code without alteration.
