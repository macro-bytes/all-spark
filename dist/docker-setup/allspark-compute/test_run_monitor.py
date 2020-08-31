import unittest
import os
from run_monitor import get_app_exit_status, get_cluster_status, APP_EXIT_STATUS_PATH

class TestRunMonitor(unittest.TestCase):

    def test_get_cluster_status(self):
        cluster_status = get_cluster_status()
        assert "ALIVE" == cluster_status["status"]
        assert [] == cluster_status["completedapps"]

    def test_get_app_exit_status(self):
        def set_exit_failure():
            with open(APP_EXIT_STATUS_PATH, "w") as fh:
                fh.write("ERROR")

        try:
            os.remove(APP_EXIT_STATUS_PATH)
        except:
            ...
        exit_status = get_app_exit_status()
        assert "" == exit_status

        set_exit_failure()
        exit_status = get_app_exit_status()
        assert "ERROR" == exit_status

if __name__ == '__main__':
    unittest.main()

