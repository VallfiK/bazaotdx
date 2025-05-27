from selenium import webdriver
from selenium.webdriver.common.by import By
from selenium.webdriver.chrome.service import Service
from selenium.webdriver.support.ui import WebDriverWait
from selenium.webdriver.support import expected_conditions as EC
from webdriver_manager.chrome import ChromeDriverManager
import time
import random
from fake_useragent import UserAgent

# Конфигурация
NUMBER_OF_REQUESTS = 20  # Количество заявок
BASE_DELAY = 2           # Базовая задержка между действиями (сек)
USER_AGENT = UserAgent().random  # Случайный User-Agent

# Настройка ChromeOptions
def get_chrome_options():
    chrome_options = webdriver.ChromeOptions()
    
    # Скрытие автоматизации
    chrome_options.add_argument("--disable-blink-features=AutomationControlled")
    chrome_options.add_experimental_option("excludeSwitches", ["enable-automation"])
    chrome_options.add_experimental_option('useAutomationExtension', False)
    
    # Обход блокировок
    chrome_options.add_argument('--no-sandbox')
    chrome_options.add_argument('--disable-dev-shm-usage')
    chrome_options.add_argument('--ignore-certificate-errors')
    chrome_options.add_argument('--allow-running-insecure-content')
    chrome_options.add_argument(f'user-agent={USER_AGENT}')
    
    # Отключение расширений и функций
    chrome_options.add_argument("--disable-extensions")
    chrome_options.add_argument("--disable-infobars")
    chrome_options.add_argument("--disable-notifications")
    
    return chrome_options

# Инициализация драйвера
def init_driver():
    service = Service(ChromeDriverManager().install())
    return webdriver.Chrome(service=service, options=get_chrome_options())

# Случайная задержка
def random_delay(min=1, max=3):
    time.sleep(random.uniform(min, max))

def human_like_click(driver, xpath):
    try:
        element = WebDriverWait(driver, 10).until(
            EC.element_to_be_clickable((By.XPATH, xpath))
        )
        driver.execute_script("arguments[0].scrollIntoView({behavior: 'smooth'});", element)
        random_delay(0.5, 1.5)
        element.click()
    except Exception as e:
        raise Exception(f"Click error: {str(e)}")

def human_like_input(driver, xpath, text):
    try:
        element = WebDriverWait(driver, 10).until(
            EC.element_to_be_clickable((By.XPATH, xpath))
        )
        random_delay(0.2, 0.7)
        element.clear()
        for char in text:
            element.send_keys(char)
            time.sleep(random.uniform(0.05, 0.2))
    except Exception as e:
        raise Exception(f"Input error: {str(e)}")

def submit_request(driver, attempt):
    try:
        # Плавный скролл вместо резкого перемещения
        driver.execute_script("window.scrollBy({top: 430, behavior: 'smooth'});")
        random_delay()
        
        # Используем более стабильные селекторы
        human_like_click(driver, "//button[contains(text(), 'Оформить')]")
        human_like_click(driver, "//button[contains(@class, 'confirm-btn')]")
        
        # Заполнение формы
        human_like_input(driver, "//input[@placeholder='Фамилия']", "Шепелевич")
        human_like_input(driver, "//input[@placeholder='Имя']", "Ольга")
        human_like_input(driver, "//input[@placeholder='Телефон']", "89508014178")
        human_like_input(driver, "//input[@placeholder='Email']", "9508014178@mail.ru")
        human_like_input(driver, "//textarea", "сдэк:г. Липецк, ул. Бехтеева С. С., д. 9")
        
        # Финальное подтверждение
        human_like_click(driver, "//button[@type='submit']")
        print(f"Заявка {attempt} успешно оформлена!")
        return True
    
    except Exception as e:
        print(f"Ошибка при оформлении заявки {attempt}: {str(e)}")
        return False

def main():
    driver = init_driver()
    try:
        driver.get("https://f-ariel.ru/auth/")
        random_delay(15, 25)  # Имитация поведения пользователя
        
        for i in range(1, NUMBER_OF_REQUESTS + 1):
            try:
                driver.get("https://order.f-ariel.ru/")
                WebDriverWait(driver, 15).until(
                    EC.presence_of_element_located((By.TAG_NAME, 'body'))
                
                if not submit_request(driver, i):
                    # Попытка восстановления сессии
                    driver.delete_all_cookies()
                    random_delay(3, 5)
                    driver.refresh()
                    
                # Случайная задержка между запросами
                random_delay(BASE_DELAY, BASE_DELAY * 2)
                
            except Exception as e:
                print(f"Критическая ошибка: {str(e)}. Перезапуск драйвера...")
                driver.quit()
                driver = init_driver()
                random_delay(5, 10)
                
    finally:
        driver.quit()

if __name__ == "__main__":
    main()