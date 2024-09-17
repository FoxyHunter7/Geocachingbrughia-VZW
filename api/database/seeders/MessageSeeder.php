<?php

namespace Database\Seeders;

use App\Models\Message;
use Illuminate\Database\Seeder;
use Illuminate\Support\Facades\DB;

class MessageSeeder extends Seeder
{
    /**
     * Run the database seeds.
     */
    public function run(): void
    {
        DB::table('messages')->truncate();

        $data = ReadFromCsv::getDataFromCsv(storage_path('app/data/csv/messages.csv'), ';');

        $model = new Message();
        foreach ($data as $row) {
            $model->create($row);
        }
    }
}
