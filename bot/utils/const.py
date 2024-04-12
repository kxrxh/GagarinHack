import os


class Constant:
    @staticmethod
    def get_yandex_url() -> str:
        return "https://llm.api.cloud.yandex.net/foundationModels/v1/completion"

    @staticmethod
    def get_gigachat_url() -> str:
        # todo: change
        return ""

    @staticmethod
    def get_bot_key() -> str:
        return os.environ['BOT_KEY']

    @staticmethod
    def get_ya_dir_id() -> str:
        return os.environ['YA_DIR_ID']

    @staticmethod
    def get_ya_api_key() -> str:
        return os.environ['YA_API_KEY']
