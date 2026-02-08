
import time
from playwright.sync_api import sync_playwright, expect

def run(playwright):
    browser = playwright.chromium.launch(headless=True)
    context = browser.new_context()
    page = context.new_page()

    page.on("console", lambda msg: print(f"Browser Console: {msg.text}"))
    page.on("pageerror", lambda exc: print(f"Browser Error: {exc}"))

    try:
        print("Navigating to Dashboard...")
        page.goto("http://localhost:3000")
        time.sleep(1) # Wait for load

        # 1. Verify System Tab exists and works
        print("Clicking System Tab...")
        # Use specific locator for the nav link to avoid ambiguity with "System Information" header
        page.get_by_role("link", name="System").click()
        time.sleep(1)

        print("Verifying System Info content...")
        expect(page.get_by_text("System Information")).to_be_visible()
        expect(page.get_by_text("Build Info")).to_be_visible()

        # Verify dynamic content loading (version shouldn't be 'Loading...')
        # It should be 1.3.1
        expect(page.locator("#sysAppVersion")).to_contain_text("1.3.1")

        page.screenshot(path="verification_system_tab.png")
        print("System Tab verified.")

        # 2. Verify Custom Client 'Config Key' field
        print("Clicking Clients Tab...")
        page.click("text=Clients")
        time.sleep(1)

        print("Verifying Config Key field...")
        expect(page.get_by_text("Config Key (Optional)")).to_be_visible()
        page.screenshot(path="verification_custom_client_enhanced.png")
        print("Custom Client UI verified.")

    except Exception as e:
        print(f"Test failed: {e}")
        page.screenshot(path="verification_failure_v2.png")
        raise e
    finally:
        browser.close()

with sync_playwright() as playwright:
    run(playwright)
