package io.temporal.omes;

public class ActivitiesImpl implements Activities {

  @Override
  public void noopActivity() {}

  @Override
  public String echo(String input) {
    return input;
  }

  @Override
  public void delay(int milliseconds) throws InterruptedException {
    Thread.sleep(milliseconds);
  }
}
