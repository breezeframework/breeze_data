#!/bin/bash
# Получаем последнюю версию в формате vX.Y.Z
latest_tag=$(git describe --tags $(git rev-list --tags --max-count=1))

# Извлекаем номера версий
major=$(echo $latest_tag | cut -d. -f1 | cut -c2-)
minor=$(echo $latest_tag | cut -d. -f2)
patch=$(echo $latest_tag | cut -d. -f3)

# Увеличиваем номер патча
new_patch=$((patch + 1))
new_version="v${major}.${minor}.${new_patch}"

# Создаем новый тег и пушим его
git tag $new_version
git push origin $new_version
