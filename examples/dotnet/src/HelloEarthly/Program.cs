using System;
using NodaTime;

namespace HelloEarthly
{
    class Program
    {
        static void Main(string[] args)
        {
            Console.WriteLine("Hello World! Its " + SystemClock.Instance.GetCurrentInstant());
        }
    }
}
