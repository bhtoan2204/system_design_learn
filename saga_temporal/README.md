This project demonstrates how to use Temporal to implement Saga Orchestration for a simple travel booking process
The workflow consists of the following steps:

booking hotel -> booking flight -> booking car

If booking the car fails, revert (cancel) the flight
If booking the flight fails, revert (cancel) the hotel

This ensures that the overall workflow maintains consistency across multiple services:
- All steps must succeed.
- If any step fails, all previous steps are rolled back to ensure data consistency.