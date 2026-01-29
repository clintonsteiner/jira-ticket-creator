"""
Example: Using JIRA Ticket Creator from Python

This example demonstrates how to use the JiraClient
to create tickets directly from Python code.
"""

from jira_client import JiraClient


def main():
    # Initialize client
    client = JiraClient(
        url="https://company.atlassian.net",
        email="user@company.com",
        token="your-api-token",
        project="PROJ"
    )

    print(f"Client initialized. Library version: {client.get_version()}")

    # Example 1: Create a simple task
    print("\n[Example 1] Creating a simple task...")
    try:
        ticket = client.create_ticket(
            summary="Fix login page bug",
            description="Login button not responding on mobile",
            issue_type="Bug",
            priority="High"
        )
        print(f"Created: {ticket['key']}")
        print(f"URL: {ticket['url']}")
    except Exception as e:
        print(f"Error: {e}")

    # Example 2: Create a story with labels
    print("\n[Example 2] Creating a story with labels...")
    try:
        ticket = client.create_ticket(
            summary="Implement OAuth 2.0 authentication",
            description="Add OAuth 2.0 support to the platform",
            issue_type="Story",
            priority="High",
            labels=["auth", "security", "oauth"],
            assignee="john@company.com"
        )
        print(f"Created: {ticket['key']}")
    except Exception as e:
        print(f"Error: {e}")

    # Example 3: Create a task with blockers
    print("\n[Example 3] Creating a task with blockers...")
    try:
        ticket = client.create_ticket(
            summary="Update API endpoints",
            issue_type="Task",
            priority="Medium",
            blocked_by=["PROJ-123", "PROJ-124"]
        )
        print(f"Created: {ticket['key']}")
    except Exception as e:
        print(f"Error: {e}")

    # Example 4: Extract project from ticket key
    print("\n[Example 4] Extracting project key from ticket...")
    try:
        project = client.extract_project_key("PROJ-123")
        print(f"Project extracted: {project}")
    except Exception as e:
        print(f"Error: {e}")

    # Example 5: Create client with ticket key
    print("\n[Example 5] Creating client from ticket key...")
    try:
        client2 = JiraClient(
            url="https://company.atlassian.net",
            email="user@company.com",
            token="your-api-token",
            ticket="PROJ-100"  # Project auto-extracted
        )
        print(f"Client created with project: {client2.project}")
    except Exception as e:
        print(f"Error: {e}")

    # Example 6: Get ticket details
    print("\n[Example 6] Getting ticket details...")
    try:
        ticket = client.get_ticket("PROJ-123")
        print(f"Ticket: {ticket['key']}")
        print(f"Summary: {ticket['summary']}")
        print(f"Status: {ticket['status']}")
    except Exception as e:
        print(f"Error: {e}")

    # Example 7: Search tickets
    print("\n[Example 7] Searching for tickets...")
    try:
        results = client.search(status="In Progress")
        print(f"Found {results['count']} tickets")
        for ticket in results.get('tickets', []):
            print(f"  - {ticket['key']}: {ticket['summary']}")
    except Exception as e:
        print(f"Error: {e}")

    # Example 8: Update ticket
    print("\n[Example 8] Updating ticket...")
    try:
        response = client.update_ticket(
            "PROJ-123",
            priority="High",
            description="Updated description"
        )
        print(f"Updated: {response['key']}")
    except Exception as e:
        print(f"Error: {e}")

    # Example 9: Advanced search with JQL
    print("\n[Example 9] Advanced JQL search...")
    try:
        results = client.search(jql='project = PROJ AND priority = "Highest" AND status != "Done"')
        print(f"Critical tickets: {results['count']}")
    except Exception as e:
        print(f"Error: {e}")


if __name__ == "__main__":
    main()
