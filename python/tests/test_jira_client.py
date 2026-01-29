"""
Unit tests for JIRA Client using mocked C library
"""

import json
import unittest
from unittest import mock
from pathlib import Path
import sys

# Add parent directory to path for imports
sys.path.insert(0, str(Path(__file__).parent.parent))

from jira_client import JiraClient


class MockCLib:
    """Mock C library for testing"""

    def __init__(self):
        self.version = "1.0.0"

    def CreateTicketJSON(self, url, email, token, project, json_data):
        """Mock CreateTicketJSON"""
        data = json.loads(json_data.decode())
        response = {
            "key": "PROJ-123",
            "id": "10000",
            "url": "https://company.atlassian.net/browse/PROJ-123",
        }
        if "error" in data:
            response["error"] = data["error"]
        return json.dumps(response).encode()

    def GetTicket(self, url, email, token, key):
        """Mock GetTicket"""
        if key.decode() == "PROJ-123":
            response = {
                "key": "PROJ-123",
                "id": "10000",
                "summary": "Test ticket",
                "description": "Test description",
                "status": "To Do",
                "issue_type": "Task",
                "priority": "Medium",
                "assignee": "user@company.com",
                "labels": ["test"],
                "url": "https://company.atlassian.net/browse/PROJ-123",
            }
        else:
            response = {"error": "Ticket not found"}
        return json.dumps(response).encode()

    def SearchTickets(self, url, email, token, jql):
        """Mock SearchTickets"""
        response = {
            "total": 1,
            "count": 1,
            "tickets": [
                {
                    "key": "PROJ-123",
                    "summary": "Test ticket",
                    "status": "To Do",
                    "issue_type": "Task",
                    "priority": "Medium",
                    "assignee": "user@company.com",
                }
            ],
        }
        return json.dumps(response).encode()

    def UpdateTicket(self, url, email, token, key, json_data):
        """Mock UpdateTicket"""
        response = {"key": key.decode(), "message": "Update request received"}
        return json.dumps(response).encode()

    def ExtractProjectKey(self, ticket_key):
        """Mock ExtractProjectKey"""
        parts = ticket_key.decode().split("-")
        return parts[0].encode() if len(parts) > 1 else b""

    def FreeMemory(self, ptr):
        """Mock FreeMemory - does nothing"""
        pass

    def Version(self):
        """Mock Version"""
        return self.version.encode()


class TestJiraClientInit(unittest.TestCase):
    """Test JiraClient initialization"""

    @mock.patch("jira_client.JiraClient._load_library")
    def test_init_with_project(self, mock_load):
        """Test initialization with project key"""
        mock_lib = MockCLib()
        mock_load.return_value = mock_lib

        client = JiraClient(
            url="https://company.atlassian.net",
            email="user@company.com",
            token="token",
            project="PROJ",
        )

        self.assertEqual(client.project, "PROJ")
        self.assertEqual(client.url, "https://company.atlassian.net")

    @mock.patch("jira_client.JiraClient._load_library")
    def test_init_with_ticket(self, mock_load):
        """Test initialization with ticket key"""
        mock_lib = MockCLib()
        mock_load.return_value = mock_lib

        with mock.patch.object(
            mock_lib, "ExtractProjectKey", return_value=b"PROJ"
        ):
            client = JiraClient(
                url="https://company.atlassian.net",
                email="user@company.com",
                token="token",
                ticket="PROJ-100",
            )

            self.assertEqual(client.project, "PROJ")

    @mock.patch("jira_client.JiraClient._load_library")
    def test_init_without_project_or_ticket(self, mock_load):
        """Test initialization fails without project or ticket"""
        mock_lib = MockCLib()
        mock_load.return_value = mock_lib

        with self.assertRaises(ValueError):
            JiraClient(
                url="https://company.atlassian.net",
                email="user@company.com",
                token="token",
            )


class TestCreateTicket(unittest.TestCase):
    """Test ticket creation"""

    def setUp(self):
        """Set up test client"""
        self.mock_lib = MockCLib()

    @mock.patch("jira_client.JiraClient._load_library")
    def test_create_simple_ticket(self, mock_load):
        """Test creating a simple ticket"""
        mock_load.return_value = self.mock_lib

        client = JiraClient(
            url="https://company.atlassian.net",
            email="user@company.com",
            token="token",
            project="PROJ",
        )

        ticket = client.create_ticket(summary="Test ticket")

        self.assertEqual(ticket["key"], "PROJ-123")
        self.assertEqual(ticket["id"], "10000")
        self.assertIn("url", ticket)

    @mock.patch("jira_client.JiraClient._load_library")
    def test_create_ticket_with_all_fields(self, mock_load):
        """Test creating a ticket with all fields"""
        mock_load.return_value = self.mock_lib

        client = JiraClient(
            url="https://company.atlassian.net",
            email="user@company.com",
            token="token",
            project="PROJ",
        )

        ticket = client.create_ticket(
            summary="Test ticket",
            description="Test description",
            issue_type="Bug",
            priority="High",
            assignee="john@company.com",
            labels=["bug", "urgent"],
            blocked_by=["PROJ-1"],
        )

        self.assertEqual(ticket["key"], "PROJ-123")

    @mock.patch("jira_client.JiraClient._load_library")
    def test_create_ticket_without_summary_fails(self, mock_load):
        """Test that creating ticket without summary fails"""
        mock_load.return_value = self.mock_lib

        client = JiraClient(
            url="https://company.atlassian.net",
            email="user@company.com",
            token="token",
            project="PROJ",
        )

        with self.assertRaises(ValueError):
            client.create_ticket(summary="")


class TestGetTicket(unittest.TestCase):
    """Test retrieving ticket details"""

    @mock.patch("jira_client.JiraClient._load_library")
    def test_get_existing_ticket(self, mock_load):
        """Test getting an existing ticket"""
        mock_lib = MockCLib()
        mock_load.return_value = mock_lib

        client = JiraClient(
            url="https://company.atlassian.net",
            email="user@company.com",
            token="token",
            project="PROJ",
        )

        ticket = client.get_ticket("PROJ-123")

        self.assertEqual(ticket["key"], "PROJ-123")
        self.assertEqual(ticket["summary"], "Test ticket")
        self.assertEqual(ticket["status"], "To Do")

    @mock.patch("jira_client.JiraClient._load_library")
    def test_get_nonexistent_ticket(self, mock_load):
        """Test getting a nonexistent ticket fails"""
        mock_lib = MockCLib()
        mock_load.return_value = mock_lib

        client = JiraClient(
            url="https://company.atlassian.net",
            email="user@company.com",
            token="token",
            project="PROJ",
        )

        with self.assertRaises(Exception):
            client.get_ticket("PROJ-999")


class TestSearchTickets(unittest.TestCase):
    """Test ticket search"""

    @mock.patch("jira_client.JiraClient._load_library")
    def test_search_by_status(self, mock_load):
        """Test searching tickets by status"""
        mock_lib = MockCLib()
        mock_load.return_value = mock_lib

        client = JiraClient(
            url="https://company.atlassian.net",
            email="user@company.com",
            token="token",
            project="PROJ",
        )

        results = client.search(status="To Do")

        self.assertEqual(results["total"], 1)
        self.assertEqual(results["count"], 1)
        self.assertEqual(len(results["tickets"]), 1)

    @mock.patch("jira_client.JiraClient._load_library")
    def test_search_by_jql(self, mock_load):
        """Test searching with JQL query"""
        mock_lib = MockCLib()
        mock_load.return_value = mock_lib

        client = JiraClient(
            url="https://company.atlassian.net",
            email="user@company.com",
            token="token",
            project="PROJ",
        )

        results = client.search(jql='project = PROJ AND status = "To Do"')

        self.assertEqual(results["count"], 1)

    @mock.patch("jira_client.JiraClient._load_library")
    def test_search_by_key(self, mock_load):
        """Test searching by ticket key"""
        mock_lib = MockCLib()
        mock_load.return_value = mock_lib

        client = JiraClient(
            url="https://company.atlassian.net",
            email="user@company.com",
            token="token",
            project="PROJ",
        )

        results = client.search(key="PROJ-123")

        self.assertEqual(results["count"], 1)

    @mock.patch("jira_client.JiraClient._load_library")
    def test_search_default_project(self, mock_load):
        """Test search defaults to project"""
        mock_lib = MockCLib()
        mock_load.return_value = mock_lib

        client = JiraClient(
            url="https://company.atlassian.net",
            email="user@company.com",
            token="token",
            project="PROJ",
        )

        results = client.search()

        self.assertEqual(results["count"], 1)


class TestUpdateTicket(unittest.TestCase):
    """Test ticket updates"""

    @mock.patch("jira_client.JiraClient._load_library")
    def test_update_ticket_summary(self, mock_load):
        """Test updating ticket summary"""
        mock_lib = MockCLib()
        mock_load.return_value = mock_lib

        client = JiraClient(
            url="https://company.atlassian.net",
            email="user@company.com",
            token="token",
            project="PROJ",
        )

        response = client.update_ticket("PROJ-123", summary="New summary")

        self.assertEqual(response["key"], "PROJ-123")

    @mock.patch("jira_client.JiraClient._load_library")
    def test_update_multiple_fields(self, mock_load):
        """Test updating multiple fields"""
        mock_lib = MockCLib()
        mock_load.return_value = mock_lib

        client = JiraClient(
            url="https://company.atlassian.net",
            email="user@company.com",
            token="token",
            project="PROJ",
        )

        response = client.update_ticket(
            "PROJ-123", priority="High", assignee="john@company.com"
        )

        self.assertEqual(response["key"], "PROJ-123")

    @mock.patch("jira_client.JiraClient._load_library")
    def test_update_without_fields_fails(self, mock_load):
        """Test that updating without fields fails"""
        mock_lib = MockCLib()
        mock_load.return_value = mock_lib

        client = JiraClient(
            url="https://company.atlassian.net",
            email="user@company.com",
            token="token",
            project="PROJ",
        )

        with self.assertRaises(ValueError):
            client.update_ticket("PROJ-123")


class TestExtractProjectKey(unittest.TestCase):
    """Test project key extraction"""

    @mock.patch("jira_client.JiraClient._load_library")
    def test_extract_valid_key(self, mock_load):
        """Test extracting project key from ticket"""
        mock_lib = MockCLib()
        mock_load.return_value = mock_lib

        client = JiraClient(
            url="https://company.atlassian.net",
            email="user@company.com",
            token="token",
            project="PROJ",
        )

        project = client.extract_project_key("PROJ-123")

        self.assertEqual(project, "PROJ")

    @mock.patch("jira_client.JiraClient._load_library")
    def test_extract_invalid_key_fails(self, mock_load):
        """Test extracting from invalid ticket key fails"""
        mock_lib = MockCLib()
        mock_load.return_value = mock_lib

        client = JiraClient(
            url="https://company.atlassian.net",
            email="user@company.com",
            token="token",
            project="PROJ",
        )

        with self.assertRaises(ValueError):
            client.extract_project_key("INVALID")


class TestGetVersion(unittest.TestCase):
    """Test version retrieval"""

    @mock.patch("jira_client.JiraClient._load_library")
    def test_get_version(self, mock_load):
        """Test getting library version"""
        mock_lib = MockCLib()
        mock_load.return_value = mock_lib

        client = JiraClient(
            url="https://company.atlassian.net",
            email="user@company.com",
            token="token",
            project="PROJ",
        )

        version = client.get_version()

        self.assertEqual(version, "1.0.0")


if __name__ == "__main__":
    unittest.main()
